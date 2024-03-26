../gravity-cli product create drink \
    --schema "./assets/schema.json" \
    --enabled

../gravity-cli product ruleset add drink drinkCreated \
    --enabled \
    --event=drinkCreated --method=create \
    --schema="./assets/schema.json"