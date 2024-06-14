Feature: Data Product subscribe 

Scenario:
Given NATS has been opened
Given Dispatcher has been opened
#Scenario
	Scenario:  Successful scenario. Use the `product sub` command to receive all data published to the specified data product.
	Given Create data product "'drink'"
    Given Create data product "'drink'" with ruleset "'drinkCreated'"
    Given Publish 10 Events
    Given Set subscribe Timeout with "'1'" seconnd
    When Subscribe data product "'<ProductName>'" using parameters "'<SubName>'" "'<Partitions>'" "'<Seq>'"
    Then The CLI returns all events data within the "'<Partitions>'" and after "'<Seq>'"
    Examples:
    |  ID  | ProductName | SubName |   Partitions   |      Seq      |
    | M(1) |   drink     |         |      -1        |       1       |
    | M(2) |   drink     |         |      0         |       1       |
    | M(3) |   drink     |         |      200       |       1       |
    | M(4) |   drink     |         |  2147483647    |       1       |
    | M(5) |   drink     |         |  -2147483647   |       1       |
    | M(6) |   drink     |         |    [ignore]    |       1       |
    | M(7) |   drink     |         |   2147483648   |       1       |
    | M(8) |   drink     |         |  -2147483648   |       1       |
    | M(9) |   drink     |         |    131,200     |       1       |
    | M(10)|   drink     |         |      -1        |   4294967295  |
    | M(11)|   drink     |         |      -1        |    [ignore]   |
    | M(12)|   drink     |         |      -1        |   4294967296  |
    | M(13)|   drink     |         |      -1        |       5       |

#Scenario
	Scenario: Failure scenario. Use the `product sub` command to receive all data published to the specified data product.
	Given Create data product "'drink'"
    Given Create data product "'drink'" with ruleset "'drinkCreated'"
    Given Publish 10 Events
    When Subscribe data product "'<ProductName>'" using parameters "'<SubName>'" "'<Partitions>'" "'<Seq>'"
    Then The CLI returns exit code 1
    Examples:
    |   ID   | ProductName | SubName    |   Partitions   |      Seq      |
    |  E1(1) |   notExist  |            |      -1        |       1       |
    |  E1(2) |     drink   |            |    notNumber   |       1       |
    |  E1(3) |     drink   |            |      0.1       |       1       |
    |  E1(4) |     drink   |            |      ,         |       1       |
    |  E1(5) |     drink   |            |    abc,200     |       1       |
    |  E1(6) |     drink   |            |      ""        |       1       |
    |  E1(7) |     drink   |            |      -1        |       0       |
    |  E1(8) |     drink   |            |      -1        |      -1       |  
    |  E1(9) |     drink   |            |      -1        |   notNumber   | 
    |  E1(10)|     drink   |            |      -1        |      0.1      |
    |  E1(11)|     drink   |            |      -1        |      ""       | 