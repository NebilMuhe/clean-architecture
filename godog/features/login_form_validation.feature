Feature: Login User

    Scenario: validate the user input
        Given User is on login page
        When User enters "<Username>" and "<Password>"
        Then The system sholud return an error "<ErrorMessage>"
        Examples:
        | Username   | Password              | ErrorMessage                                       |
        |            | 12QWas@#              | username required                                  |
        | Matheo     |                       | password required                                  |
        | 123455     | 12QWas@#              | username must be valid                            |
        | dave       | 12QWas@#              | Username length must be atleast 5 characters       |
        | david      | 12QWas                | Password length must be atleast 8 characters long  |
        | david      | 12345678              | Password must contain atleast one uppercase letters,one lowercase letters, digits and special characters  |
        | david      | 1234ABCD              | Password must contain atleast one uppercase letters,one lowercase letters, digits and special characters  |
        | david      | 12ABCDab              | Password must contain atleast one uppercase letters,one lowercase letters, digits and special characters  |
        | david1     | 12ABCD%$              | Password must contain atleast one uppercase letters,one lowercase letters, digits and special characters  |
        