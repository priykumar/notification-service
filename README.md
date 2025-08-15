Problem Statement:
A notification service that facilitates the sending of notifications thorugh email, slack, in-app notification 

Requirements:
1. Send Notifications through channels like Email, Slack, In-app
2. Producers tells notification service
    a. to send the message instantly
    b. to send the message in future
3. Producers can design their own template

APIs
1. POST /producer/template
Request Body:
{
    "name":
    "subject":
    "content":
}

2. POST /consumer/register&type=[slack/inapp/email]
Request Body:
{
    "detail": "slackId/mobileId/EmailId"
}

Data model
1. Notification:
    a. to
    b. from
    c. template
    d. time
    e. content


Payment Reminder
"Dear {0}, please clear your pending payment before {1}. Thank you."

Job Recommendation
"Hi {0}, we found a job matching your skills. Apply soon!"

Product Sale
"Exclusive offer! Get {0} at discounted price until {1}. Hurry!"

curl -X POST http://localhost:8080/producer/template \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Payment Reminder",
    "subject": "Payment Due Reminder",
    "message": "Dear {0}, please clear your pending payment before {1}. Thank you."
  }'



