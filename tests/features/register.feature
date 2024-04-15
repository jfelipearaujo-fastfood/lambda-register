Feature: register
    In order to register a new user
    As a customer
    I need to be able to register with my CPF and password

    Scenario: Register a non anonymous user
        Given the user CPF is "548.644.620-97"
        And the user password is "12345678"
        When the user request to be registered
        Then the user should be registered successfully

    Scenario: Register an anonymous user
        Given the user CPF is ""
        And the user password is ""
        When the user request to be registered
        Then the user should be registered successfully