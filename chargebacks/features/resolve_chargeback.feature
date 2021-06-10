Feature: resolve chargeback
  In order to resolve a chargeback
  As an API user
  I need to be able

  Scenario: does not allow solve a chargeback with answer invalid
    When i send PUT request to "resolveChargeback"
    Then the response code should be 400
