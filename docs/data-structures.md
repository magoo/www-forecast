## Data Structures

### Entity Identifier
+ id (number)


### User
There is no explicit user model currently. There is a strict, external dependency on Google SSO to assign entities to a user.


### Question
A common set of fields that are composed into other sub-types such as an `Estimate`.

#### Properties
+ Include Entity Identifier
+ ownerid - owner of the question aka moderator
+ date - RFC 3339 UTC creation date
+ hd - hosted domain (google apps); owning organization e.g. acme.com or public
+ title - of question
+ brierscore (number) - squared error (60% sure of something means a 40%^2 error...averaged over time aka historical accuracy at predicting something); like a reputation score (could be used for users, groups, questions, topics, etc)
+ concluded (boolean) - is this question done?
+ concludedtime - RFC 3339 UTC conclusion date
+ records (array[string]) - audit history; snapshot of past states of the question
+ url - UUID for public sharing e.g. a permalink
+ type - specifies a human readable label for the question's type e.g. "Scenario"


### Answer
A common set of fields that are composed into other sub-types such as a `Range`.

#### Properties
+ Include Entity Identifier
+ ownerid - owner of the answer aka moderator
+ hd - hosted domain (google apps); owning oranization e.g. acme.com or public
+ date - RFC 3339 UTC creation date
+ useralias - fake name of participant for anonymizing users so not to affect other user's answers
+ url - permalink for public sharing (not currently used for `Answer`)
+ title - apparently not used
+ description - apparently not used
+ brierscore (number) - final score of this answer based on actual outcome
+ concluded (boolean) - whether or not this answer's question has been concluded




### Estimate
A kind of `Question` that attempts to find the range in which an outcome will be within. `Estimate` may be too vague.

#### Properties
+ Include Question
+ avgminimum - based on given answers
+ avgmaximum - based on given answers
+ actual - hidden value, revealed when question is concluded
+ unit: jelly beans - what's being estimated


### Range (rename to interval)
confidence level is static, set to 95% for all answers for now
+ Included Answer
+ minimum: 100 (number)
+ maximum: 300 (number)




### Scenario
A `Question` that attempts to determine which outcome is most likely to be the actual outcome via a percentage. Would this be better named as a `Probability`?

#### Properties
+ Include Question
+ options: yes, no (array[string]) - list of possible choices
+ results: 30, 70 (array[number]) - should be a list of float values for each option (assumes sequential/consistent order). Value is empty at first. If 3 people say yes and 7 say no, then its [30, 70] directly corresponding to the order of the given `options`.
+ resultsindex: 1 (number) - which result from the results list ended up being true


### Forecasts

#### Properties
+ Includes Answer
+ forecasts: 60.0, 40.0 (array[number]) - list of given percent values for each option e.g. 60% likely to be yes, and 40% no




### Rank
A `Question` that attempts to determine the actual order of a set of possible outcomes. Results are calculated on the fly.

#### Properties
+ Include Question
+ options: A, B, C (array[string]) - list of options


### Sort

#### Properties
+ Include Answer
+ options: 2, 1, 3 (array[number]) - ordered list of indexes corresponding the the list of possible options
