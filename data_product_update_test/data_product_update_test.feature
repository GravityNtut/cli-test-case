Feature: Data Product update

Scenario:
Given NATS has been opened
Given Dispatcher has been opened
#Scenario
	Scenario Outline: Success scenario for updating a data product
	Given Create data product "'drink'" and enabled is "'<GivenDPEnabled>'"
	When Update the name of data product to "'<ProductName>'" using parameters "'<Description>'" "'<Enabled>'" "'<Schema>'"
	Then Check updating data product success
	And Use NATS jetstream to query the "'drink'" update successfully and parameters are "'<Description>'" "'<Schema>'" "'<Enabled>'" 
	Examples:
	|  ID  | ProductName | Description  |		 Schema     	| Enabled   | GivenDPEnabled |
	| M(1) | drink       | [ignore]     |		[ignore]   		| [ignore]  |   [false]    |
	| M(2) | drink       | description  |		[ignore] 		| [ignore]  |   [false]    |
	| M(3) | drink       | ""  			|		[ignore]   		| [ignore]  |   [false]    |
	| M(4) | drink       | " "    	    |	    [ignore]   	    | [ignore]  |   [false]    |
	| M(5) | drink       | [ignore]     | ./assets/schema.json  | [ignore]  |   [false]    |
	| M(6) | drink       | [ignore]     |		[ignore]		| [true]    |   [false]    |
	| M(7) | drink       | [ignore]     |		[ignore]		| [false]   |   [true]     |
	| M(8) | drink       |  description | ./assets/schema.json  | [true]    |   [false]    |

#Scenario
	Scenario Outline: Fail scenario for updating data product
	Given Create data product "'drink'" and enabled is "'[true]'"
	Given Store NATS copy of existing data product "'drink'"
	When Update the name of data product to "'<ProductName>'" using parameters "'<Description>'" "'<Enabled>'" "'<Schema>'"
	Then CLI returns exit code 1
	# And The error message should be "'<Error_message>'"
	And Use NATS jetstream to query the "'drink'" without changing parameters
	Examples:
	|  ID   | ProductName | Description  | 	    Schema         		 | Enabled   | Error_message   |
	| E1(1) | not_exist   | [ignore]  	 | 		[ignore] 			 | [false]  |                 |
	| E1(2) | [null]      | [ignore] 	 |		[ignore]  		     | [false]  | 			       |
	| E1(3) | drink		  | [ignore]     | 		""					 | [false]  |                 |
	| E1(4) | drink		  | [ignore]     | ./assets/fail_schema.json | [false]  |                 |
	| E1(5) | drink		  | [ignore]     |		not_exist.json		 | [false]  |                 |