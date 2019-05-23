Feature: Moderator, Conclude Question

	Rule: The Moderator created the question.

	Background:
	As a Moderator
	I can conclude a question
	So that I can receive the Brier Score of the answer, and so panelists can observe their Brier Scores.

	Scenario: Moderator considers concluding question
		Given that I am a moderator
		When I click on the "conclude" button on a question
		Then a dialogue will appear
		And it will describe what conclusion means
		And if there is an outcome to the question
		Then it will show the optional outcomes to the Moderator
		And show the option to escape the dialogue

	Scenario: Moderator concludes the question
		Given I am a Moderator considering concluding a question
		And I have selected the outcome
		And I have selected "Conclude"
		Then all of the panelist answers will be checked for accuracy
		And panelist Brier Scores will be stored with their
		And an average Brier Score will be stored with the question

	Scenario: Moderator escapes the conclusion dialogue
		Given I am a Moderator considering concluding a question
		When I click to close or escape the dialogue
		Then I am returned to the Question without any changes
