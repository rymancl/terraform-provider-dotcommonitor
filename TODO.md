# TODO
- ~~Add license headers~~
- ~~Determine allowed task request_type's (GET,POST,PUT ??)~~
- Prevent state from being written when API errors - this might be fine now?
- ~~Update README to be proper for a TF provider~~
- ~~Make a note in README that groups must be assigned via the device, not the other way around~~

- Schema Validation - [Example validators](https://github.com/terraform-providers/terraform-provider-aws/blob/9ac9c21243ebc01ac987a85ec30d8f8dfed8170b/aws/validators.go)
  - ~~String: `ValidateFunc: validation.StringLenBetween(1, 255),`~~
  - ~~Int: `ValidateFunc: validation.IntAtLeast(0),`~~
  - ~~Task: Add IP range validation to task host IP~~
  - ~~Group: Better validation for Addresses number (valid int and valid length)~~
  
- Tasks
  - ~~Task: Add IP range validation to task host IP~~
  - ~~Task: Support for custom DNS hosts under DNS Options~~
  - Task: Rework keyword1,keyword2,keyword3 into a list with max length = 3 ?
  - ~~Task: SSL / Certificate Check (check API)~~

- All
  - Tests
  - See if I can use TypeSet instead of TypeList in schema
  - ~~Add logic to resources to check if resource exists before any CRUD operation~~
  - Should we enforce unique names from the provider side? (query API to see if the name already exists)

- Filters
  - Implement filter resource 
  - Implement filter data source
  
- Scheduler
  - Implement schedule resource 
  - Implement schedule data source