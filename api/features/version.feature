Feature: get version
  In order to know godog version
  As an API user
  I need to be able to request version

  Scenario: does not allow POST method
    Given a GET method for get version
    When i send "POST" request to "/version"
    Then the response code should be 405
    And the response should match json:
      """
      {
        "error": "Method not allowed"
      }
      """

  Scenario: should get version number
    Given a GET method for get version
    When i send "GET" request to "/version"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "version": "v0.11.0"
      }
      """