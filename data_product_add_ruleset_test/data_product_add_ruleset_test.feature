Feature: Data Product ruleset add

Scenario:
    Given NATS has been opened
    Given Dispatcher has been opened

#Scenario
    Scenario Outline: Success scenario for adding ruleset to data product
    Given Create data product "'drink'" and enabled is "'[true]'"
    When "'<ProductName>'" add ruleset "'<Ruleset>'" using parameters "'<Method>'" "'<Event>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'"
    Then Check adding ruleset success
    Then Use gravity-cli to query the "'<ProductName>'" "'<Ruleset>'" exists
    Then Use NATS jetstream to query the "'<ProductName>'" "'<Ruleset>'" exists and parameters "'<Method>'" "'<Event>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'" are correct
    Examples:
        |  ID   | ProductName | Ruleset       | Method       | Event         |   Pk          |  Desc             |          Handler_script       |           Schema               | Enabled |
        | M(1)  | drink       | drinkCreated  | create       | drinkCreated  |   id          |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(2)  | drink       | _-*=_?+@      | create       | drinkCreated  |   id          |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(3)  | drink       | 中文          |  create      | drinkCreated  |   id          |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(4)  | drink       | [a]x32768     | create       | drinkCreated  |   id          |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(5)  | drink       | drinkCreated  | " "          | " "           |   " "         |  " "              |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(6)  | drink       | drinkCreated  | _-*=_?+@     | _-*=_?+@      |   _-*=_?+@    |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(7)  | drink       | drinkCreated  | 中文         | 中文           |     中文      |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(8)  | drink       | drinkCreated  | [a]x32768    | [a]x32768     |   id          |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(9)  | drink       | drinkCreated  | "create"     | "drinkCreated"|   "id"        | "drink data desc" |     "./assets/handler.js"     |      "./assets/schema.json"    | [true]  |
        | M(10) | drink       | drinkCreated  | create       | drinkCreated  |      id,id2   |  description      |     ./assets/handler.js       |      "./assets/schema.json"    | [true]  |
        | M(11) | drink       | drinkCreated  | create       | drinkCreated  |   ""          |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(12) | drink       | drinkCreated  | create       | drinkCreated  |   [ignore]    |    [ignore]       |           [ignore]            |            [ignore]            | [true]  |
        | M(13) | drink       | drinkCreated  | create       | drinkCreated  |   [a]x32768   |   [a]x32768       |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(14) | drink       | drinkCreated  | create       | drinkCreated  |   "id, "      |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(15) | drink       | drinkCreated  | create       | drinkCreated  |   id          |  ""               |     ./assets/handler.js       |      ./assets/schema.json      | [true]  |
        | M(16) | drink       | drinkCreated  | create       | drinkCreated  |   id          |  description      |     ./assets/new.txt          |      ./assets/schema.json      | [true]  |
        | M(17) | drink       | drinkCreated  | create       | drinkCreated  |   id          |  description      |     ./assets/fail_handler.js  |      ./assets/schema.json      | [true]  |
        | M(18) | drink       | drinkCreated  | create       | drinkCreated  |   id          |  description      |     ./assets/WrongFile.jpg    |      ./assets/schema.json      | [true]  |
        | M(19) | drink       | drinkCreated  | create       | drinkCreated  |   id          |  description      |     ./assets/empty.js         |      ./assets/schema.json      | [true]  |
        | M(20) | drink       | drinkCreated  | create       | drinkCreated  |   id          |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [false] |
        | M(21) | drink       | drinkCreated  | create       | drinkCreated  |   id          |  description      |     ./assets/handler.js       |      ./assets/schema.json      | [ignore]|


#Scenario
    Scenario Outline: Scenario for repeatedly adding ruleset to data product
    Given Create data product "'drink'" and enabled is "'[true]'"
    When "'<ProductName>'" add ruleset "'<Ruleset1>'" using parameters "'<Method>'" "'<Event1>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'"
    Then Check adding ruleset success
    Then Use gravity-cli to query the "'<ProductName>'" "'<Ruleset1>'" exists
    Then Use NATS jetstream to query the "'<ProductName>'" "'<Ruleset1>'" exists and parameters "'<Method>'" "'<Event1>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'" are correct
    When "'<ProductName>'" add ruleset "'<Ruleset2>'" using parameters "'<Method>'" "'<Event2>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'"
    Then CLI returns exit code 1
    # And The error message should be "'<Error_message>'"
    Examples:
        | ID   | ProductName | Ruleset1      | Ruleset2       | Method  | Event1        | Event2         |   Pk  | Desc             |   Handler_script    |       Schema         | Enabled |  Error_message             |
        | E1(1)| drink       | drinkCreated  | drinkCreated   | create  | drinkCreated  | drinkCreated2  |   id  |  "description"   | ./assets/handler.js | ./assets/schema.json | [true]  |                            |
        | E1(2)| drink       | drinkCreated  | drinkCreated2  | create  | drinkCreated  | drinkCreated   |   id  |  "description"   | ./assets/handler.js | ./assets/schema.json | [true]  |                            |
        				
#Scenario
    Scenario Outline: Fail scenario for adding ruleset to data product
    Given Create data product "'drink'" and enabled is "'[true]'"
    When "'<ProductName>'" add ruleset "'<Ruleset>'" using parameters "'<Method>'" "'<Event>'" "'<Pk>'" "'<Desc>'" "'<Handler_script>'" "'<Schema>'" "'<Enabled>'"
    Then CLI returns exit code 1
    # And The error message should be "'<Error_message>'"
    Examples:
        | ID     | ProductName | Ruleset       | Method       | Event         |   Pk     |  Desc            |        Handler_script       |              Schema               | Enabled | Error_message |
        | E2(1)  | NotExists   |  drinkCreated |  create      | drinkCreated  |   id     |   description    |     ./assets/handler.js     |      ./assets/schema.json         | [true]  |               |
        | E2(2)  |  [null]     |    [null]     | [ignore]     |   [ignore]    | [ignore] |   [ignore]       |           [ignore]          |           [ignore]                | [true]  |               |
        | E2(3)  | drink       |    [null]     | [ignore]     |   [ignore]    | [ignore] |   [ignore]       |           [ignore]          |           [ignore]                | [true]  |               |
        | E2(4)  | drink       | drinkCreated  | ""           | ""            |   ""     |  ""              |     ./assets/handler.js     |      ./assets/schema.json         | [true]  |               |
        | E2(5)  | drink       | drinkCreated  | 中文         |  [ignore]     |   id     |   description    |     ./assets/handler.js     |      ./assets/schema.json         | [true]  |                |
        | E2(6)  | drink       | drinkCreated  | [ignore]     | drinkCreated  |   id     |   description    |     ./assets/handler.js     |     ./assets/schema.json          | [true]  |               |
        | E2(7)  | drink       | drinkCreated  | [ignore]     |  [ignore]     |   id     |   description    |     ./assets/handler.js     |     ./assets/schema.json          | [true]  |               |
        | E2(8)  | drink       | drinkCreated  | create       |  [ignore]     |   id     |   description    |     ./assets/handler.js     |     ./assets/schema.json          | [true]  |               |
        | E2(9)  | drink       | drinkCreated  | create       | drinkCreated  |   id     |   description    |     ./assets/not_exist.js   |      ./assets/schema.json         | [true]  |               |  
        | E2(10) | drink       | drinkCreated  | create       | drinkCreated  |   id     |   description    |                 ""          |      ./assets/schema.json         | [true]  |               |  
        | E2(11) | drink       | drinkCreated  | create       | drinkCreated  |   id     |   description    |     ""                      |      ""                           | [true]  |               |
        | E2(12) | drink       | drinkCreated  | create       | drinkCreated  |   id     |   description    |     ./assets/BigFile.js     |            schema.json            | [true]  |               |
        | E2(13) | drink       | drinkCreated  | create       | drinkCreated  |   id     |   description    |     ./assets/handler.js     |            not_exist.json         | [true]  |               |  
        | E2(14) | drink       | drinkCreated  | create       | drinkCreated  |   id     |   description    |     ./assets/handler.js     |      ./assets/fail_schema.json    | [true]  |               |