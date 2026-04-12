#!/usr/bin/with-contenv bashio

# Export configured environment variables for the container
export OPENEVT_ADDR="$(bashio::config 'address')"
export OPENEVT_SERIALNUMBER="$(bashio::config 'serial_number')"

echo "Starting OpenEVT with address: $OPENEVT_ADDR serial: $OPENEVT_SERIALNUMBER"

exec /usr/bin/openevt 
