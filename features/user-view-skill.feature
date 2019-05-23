Feature: User, View Skill
	To know their brier score
	A user
	Can view their skill overview page

	Scenario: Viewing the skill page
		Given a user is authenticated
		When they click on "Skill"
		Then they are directed to a page that describes their Brier Score
		And their closed forecasts
		And their current calibration
