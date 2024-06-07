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
        
    Scenario: Not registered user
        Given User is on login page
        When I send "POST" request to "/api/v1/users/login" with payload:
        """
        {
            "username": "abebe",
            "password": "12ABCD%$ab"
        }   
        """
        Then The system sholud return an error "user does not exist"
    
    Scenario: Invalid Username or Password
        Given User is on registration page
        When I send "POST" request to "/api/v1/users/register" with payload:
        """
            {
                "username": "abebe","email": "abebe@gmail.com","password": "12ABCD%$ab"
            }
        """
        And I send "POST" request to "/api/v1/users/login" with payload:
        """
            {
                "username": "abebe","password": "12ABCD#c&d"
            }
        """
        Then The system sholud return an error "invalid input"