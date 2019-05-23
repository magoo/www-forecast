Feature: User, View Skill Help
	To understand their skill and brier score
	A User
	Can view skill overview help info

	Scenario: User clicks on "Learn More"
		When a user clicks on "Learn More" from the skills page
		Then a dialogue appears with information about Brier Scores and Calibration
