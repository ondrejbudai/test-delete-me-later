#!/bin/bash
set -euo pipefail

OSBUILD_COMPOSER_TEST_DATA=/usr/share/tests/osbuild-composer

source /etc/os-release

# Colorful output.
function greenprint {
    echo -e "\033[1;32m${1}\033[0m"
}

# Provision the software under tet.
/usr/libexec/osbuild-composer-test/provision.sh

# Apply lorax patch to work around pytoml issues in RHEL 8.x.
# See BZ 1843704 or https://github.com/weldr/lorax/pull/1030 for more details.
if [[ $ID == rhel ]]; then
    sudo sed -r -i 's#toml.load\(args\[3\]\)#toml.load(open(args[3]))#' \
        /usr/lib/python3.6/site-packages/composer/cli/compose.py
    sudo rm -f /usr/lib/python3.6/site-packages/composer/cli/compose.pyc
fi

GOVC_CMD=/tmp/govc

# shellcheck source=/dev/null
source "$VCENTER_CREDS"

# We need govc to talk to vSphere
if ! hash govc; then
    greenprint "Installing govc"
    pushd /tmp
        curl -Ls --retry 5 --output govc.gz \
            https://github.com/vmware/govmomi/releases/download/v0.24.0/govc_linux_amd64.gz
        gunzip -f govc.gz
        chmod +x /tmp/govc
        $GOVC_CMD version
    popd
fi

TEST_UUID=$(uuidgen)
IMAGE_KEY=osbuild-composer-vmware-test-${TEST_UUID}

# Jenkins sets WORKSPACE to the job workspace, but if this script runs
# outside of Jenkins, we can set up a temporary directory instead.
if [[ ${WORKSPACE:-empty} == empty ]]; then
    WORKSPACE=$(mktemp -d)
fi

# Set up temporary files
TEMPDIR=$(mktemp -d)
VMWARE_CONFIG=${TEMPDIR}/vmware.toml
BLUEPRINT_FILE=${TEMPDIR}/blueprint.toml
COMPOSE_START=${TEMPDIR}/compose-start-${IMAGE_KEY}.json
COMPOSE_INFO=${TEMPDIR}/compose-info-${IMAGE_KEY}.json

# Check that the system started and is running correctly
running_test_check () {
    STATUS=$(sudo ssh -i $OSBUILD_COMPOSER_TEST_DATA/keyring/id_rsa redhat@"${1}" 'systemctl --wait is-system-running')
    if [[ $STATUS == running || $STATUS == degraded ]]; then
        echo 0
    else
        echo 1
    fi
}

# Get the compose log.
get_compose_log () {
    COMPOSE_ID=$1
    LOG_FILE=${WORKSPACE}/osbuild-${ID}-${VERSION_ID}-vmware.log

    # Download the logs.
    sudo composer-cli compose log "$COMPOSE_ID" | tee "$LOG_FILE" > /dev/null
}

# Get the compose metadata.
get_compose_metadata () {
    COMPOSE_ID=$1
    METADATA_FILE=${WORKSPACE}/osbuild-${ID}-${VERSION_ID}-vmware.json

    # Download the metadata.
    sudo composer-cli compose metadata "$COMPOSE_ID" > /dev/null

    # Find the tarball and extract it.
    TARBALL=$(basename "$(find . -maxdepth 1 -type f -name "*-metadata.tar")")
    tar -xf "$TARBALL"
    rm -f "$TARBALL"

    # Move the JSON file into place.
    cat "${COMPOSE_ID}".json | jq -M '.' | tee "$METADATA_FILE" > /dev/null
}

# Write an VMWare TOML file
tee "$VMWARE_CONFIG" > /dev/null << EOF
provider = "vmware"

[settings]
host = "${GOVMOMI_URL}"
username = "${GOVMOMI_USERNAME}"
password = "${GOVMOMI_PASSWORD}"
cluster = "${GOVMOMI_CLUSTER}"
dataStore = "${GOVMOMI_DATASTORE}"
dataCenter = "${GOVMOMI_DATACENTER}"
EOF

# Write a basic blueprint for our image.
tee "$BLUEPRINT_FILE" > /dev/null << EOF
name = "bash"
description = "A base system with bash"
version = "0.0.1"

[[packages]]
name = "bash"

[customizations.services]
enabled = ["sshd"]

[[customizations.user]]
name = "redhat"
key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC61wMCjOSHwbVb4VfVyl5sn497qW4PsdQ7Ty7aD6wDNZ/QjjULkDV/yW5WjDlDQ7UqFH0Sr7vywjqDizUAqK7zM5FsUKsUXWHWwg/ehKg8j9xKcMv11AkFoUoujtfAujnKODkk58XSA9whPr7qcw3vPrmog680pnMSzf9LC7J6kXfs6lkoKfBh9VnlxusCrw2yg0qI1fHAZBLPx7mW6+me71QZsS6sVz8v8KXyrXsKTdnF50FjzHcK9HXDBtSJS5wA3fkcRYymJe0o6WMWNdgSRVpoSiWaHHmFgdMUJaYoCfhXzyl7LtNb3Q+Sveg+tJK7JaRXBLMUllOlJ6ll5Hod root@localhost"
EOF

# Prepare the blueprint for the compose.
greenprint "???? Preparing blueprint"
sudo composer-cli blueprints push "$BLUEPRINT_FILE"
sudo composer-cli blueprints depsolve bash

# Get worker unit file so we can watch the journal.
WORKER_UNIT=$(sudo systemctl list-units | grep -o -E "osbuild.*worker.*\.service")
sudo journalctl -af -n 1 -u "${WORKER_UNIT}" &
WORKER_JOURNAL_PID=$!

# Start the compose and upload to VMWare.
greenprint "???? Starting compose"
sudo composer-cli --json compose start bash vmdk "$IMAGE_KEY" "$VMWARE_CONFIG" | tee "$COMPOSE_START"
COMPOSE_ID=$(jq -r '.build_id' "$COMPOSE_START")

# Wait for the compose to finish.
greenprint "??? Waiting for compose to finish: ${COMPOSE_ID}"
while true; do
    sudo composer-cli --json compose info "${COMPOSE_ID}" | tee "$COMPOSE_INFO" > /dev/null
    COMPOSE_STATUS=$(jq -r '.queue_status' "$COMPOSE_INFO")

    # Is the compose finished?
    if [[ $COMPOSE_STATUS != RUNNING ]] && [[ $COMPOSE_STATUS != WAITING ]]; then
        break
    fi

    # Wait 30 seconds and try again.
    sleep 30

done

# Capture the compose logs from osbuild.
greenprint "???? Getting compose log and metadata"
get_compose_log "$COMPOSE_ID"
get_compose_metadata "$COMPOSE_ID"

# Did the compose finish with success?
if [[ $COMPOSE_STATUS != FINISHED ]]; then
    echo "Something went wrong with the compose. ????"
    exit 1
fi

# Stop watching the worker journal.
sudo kill ${WORKER_JOURNAL_PID}

greenprint "???????? Building VM in vSphere"
$GOVC_CMD vm.create -u "${GOVMOMI_USERNAME}":"${GOVMOMI_PASSWORD}"@"${GOVMOMI_URL}" \
    -k=true \
    -pool="${GOVMOMI_CLUSTER}"/Resources \
    -dc="${GOVMOMI_DATACENTER}" \
    -ds="${GOVMOMI_DATASTORE}" \
    -folder="${GOVMOMI_FOLDER}" \
    -net="${GOVMOMI_NETWORK}" \
    -m=2048 -g=rhel8_64Guest -on=true -firmware=bios \
    -disk="${IMAGE_KEY}"/"${IMAGE_KEY}".vmdk \
    --disk.controller=ide \
    "${IMAGE_KEY}"

greenprint "Getting IP of created VM"
VM_IP=$($GOVC_CMD vm.ip -u "${GOVMOMI_USERNAME}":"${GOVMOMI_PASSWORD}"@"${GOVMOMI_URL}" -k=true "${IMAGE_KEY}")

# Wait for the node to come online.
greenprint "??? Waiting for VM to respond to ssh"
LOOP_COUNTER=1
while [ $LOOP_COUNTER -le 30 ]; do
    if ssh-keyscan "$VM_IP" > /dev/null 2>&1; then
        echo "SSH is up!"
        ssh-keyscan "$VM_IP" | sudo tee -a /root/.ssh/known_hosts
        break
    fi

    # ssh-keyscan has a 5 second timeout by default, so the pause per loop
    # is 10 seconds when you include the following `sleep`.
    echo "Retrying in 5 seconds..."
    sleep 5

    ((LOOP_COUNTER++))
done

greenprint "???? Checking that system is running"
for LOOP_COUNTER in {0..10}; do
    RESULT="$(running_test_check "$VM_IP")"
    if [[ $RESULT == 0 ]]; then
        echo "System is running! ????"
        break
    fi
    sleep 5
done

# Clean up
greenprint "???? Cleaning up"
$GOVC_CMD vm.destroy -u "${GOVMOMI_USERNAME}":"${GOVMOMI_PASSWORD}"@"${GOVMOMI_URL}" -k=true "${IMAGE_KEY}"

if [[ $RESULT != 0 ]]; then
    greenprint "??? Failed"
    exit 1
fi

greenprint "???? Success"
exit 0
