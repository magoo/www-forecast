Feature: User, logout
	In order to leave my session
	As a user of the application
	I want to remove my session from the browser

	Scenario:
		When the guest clicks "Sign Out"
		Then a modal opens with a confirmation dialogue
		When the guest clicks "Sign Out"
		Then their authenticated session is removed from the browser
		But if they click "x"
		Then the modal closes
		And the session is unchanged
