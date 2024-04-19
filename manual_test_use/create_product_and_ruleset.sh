../gravity-cli product create drink \
    --schema "./assets/schema.json" \
    --enabled

../gravity-cli product ruleset add drink drinkCreated \
    --enabled \
    --event=drinkCreated --method=create \
    --schema="./assets/schema.json"

../gravity-cli pub drinkCreated {"id":1,"name":"test","price":100,"kcal":50}
../gravity-cli pub drinkCreated {"id":2,"name":"Hi","price":200,"kcal":10}
../gravity-cli pub drinkCreated {"id":3,"name":"abc","price":300,"kcal":150}
../gravity-cli pub drinkCreated {"id":4,"name":"sa","price":400,"kcal":1000}
../gravity-cli pub drinkCreated {"id":5,"name":"sssa","price":500,"kcal":600}