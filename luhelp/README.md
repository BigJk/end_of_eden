# luhelp

This is the lua helper layer that exposes functions to
- Map go values to lua ones
- Map lua values to go ones
- Bind lua functions to go functions

## Naming Conventions

To conform to the naming conventions of both languages, go structs and map keys are converted to ``snake_case`` and lua table keys will be converted to ``CamelCase`` when passing data.