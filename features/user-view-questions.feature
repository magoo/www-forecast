Feature: User, Measurements
	To manage their questions
	A User
	Can view the questions they've asked.

	Scenario: The user has questions
		When the user clicks "Measurements"
		Then they are shown a list of their questions
		And each concluded question has a Brier Score

	Scenario: Not Logged In
		When a guest clicks "Measurements"
		Then they are directed to the landing page
		And shown an error

	Scenario: No Questions
		Given a user has no questions
		When a user clicks "Measurements"
		Then they are shown a dialogue that reminds them they have no questions
