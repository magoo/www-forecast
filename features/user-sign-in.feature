Feature: Guest, Login
	In order to use features
	As a guest of the application
	I want to authenticate myself and get a session

	Scenario: Successful login
		Given the guest uses a supported third party authentication
		When the guest clicks "Sign in"
		Then a modal opens with authentication options
		When the guest selects an authentication option
		Then the guest is presented with a third party authentication method
		When they authenticate with the third party
		Then they are returned to the application with a session
		But if they do not succeed
		Then they are shown an error
		And not given a session
