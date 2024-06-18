Feature: Data Product create

Scenario:
Given NATS has been opened
Given Dispatcher has been opened
#Scenario
	Scenario Outline: Success scenario for creating a data product
	When Create data product "'<ProductName>'" using parameters "'<Description>'" "'<Schema>'" "'<Enabled>'"
	Then Cli returns "'<ProductName>'" created successfully
	Then Using gravity-cli to query "'<ProductName>'" shows it exist
	Then Using NATS Jetstream to query "'<ProductName>'" shows it exist
	Examples:
	|  ID  | ProductName | Description  			| 		Schema         | Enabled |
	| M(1) | drink       | description				| ./assets/schema.json | [true]  |
	| M(2) |[a]x240      | description				| ./assets/schema.json | [true]  |
	| M(3) | drink       |     ""       		    | ./assets/schema.json | [true]  |
	| M(4) | drink       |     " " 	    		    | ./assets/schema.json | [true]  |
	| M(5) | drink       | [ignore] 			    | ./assets/schema.json | [true]  |
	| M(6) | drink       | [a]x32768    	  	    | ./assets/schema.json | [true]  |
	| M(7) | drink       | "drink data description" |"./assets/schema.json"| [true]  |
	| M(8) | drink       | description	 			|   	[ignore]	   | [true]  |
	| M(9) | drink       | description	 			| ./assets/schema.json | [false] |
	| M(10)| drink       | description	 			| ./assets/schema.json | [ignore]|


#Scenario
	Scenario Outline: Create two data products with the same name
	When Create data product "'<ProductName>'" using parameters "'<Description>'" "'<Schema>'" "'<Enabled>'"
	Then Cli returns "'<ProductName>'" created successfully
	Then Using gravity-cli to query "'<ProductName>'" shows it exist
	Then Using NATS Jetstream to query "'<ProductName>'" shows it exist
	When Create data product "'<ProductName>'" using parameters "'<Description>'" "'<Schema>'" "'<Enabled>'"
    Then CLI returns exit code 1
	# And The error message should be "'<Error_message>'"
	Examples:
	|   ID  | ProductName | Description | 		Schema         | Enabled | Error_message  |
	| E1(1) | drink       | description | ./assets/schema.json | [true]  |			    |

#Scenario
	Scenario Outline: Fail scenario for creating a data product
	When Create data product "'<ProductName>'" using parameters "'<Description>'" "'<Schema>'" "'<Enabled>'"
    Then CLI returns exit code 1
	# And The error message should be "'<Error_message>'"
	Examples:
	|   ID  | ProductName   | Description  | 			Schema       			| Enabled | Error_message |
	| E2(1) | _-*\($\)?@    | description  | 		./assets/schema.json 		| [true]  | 			  |
	| E2(2) | 中文		 	 | description  | 	./assets/schema.json  			 | [true]  |               |
	| E2(3) | [null]        | description  |	    ./assets/schema.json  		| [true]  |				  |
	| E2(4) | drink         | description  | 		notExist.json 				| [true]  |				  |
	| E2(5) | drink         | description  | 		./assets/fail_schema.json   | [true]  |				  |
	| E2(6) | drink         | description  | 			""           			| [true]  |			      |
	
