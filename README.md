# hue-config

`hue-config` allows the creation of custom scripts to control your Hue lights. You can then add these to the Alexa skill as custom 'devices' (see [Sunrise](./animations/sunrise/README.md) as an example). All hosted on AWS as lambda functions, defined and deployed using AWS SAM.

## Supported Animations

- [Sunrise](./animations/sunrise/README.md)

## Setup

The initial setup process might be a little time-consuming -- particularly if you do not have prior experience with
AWS -- but afterwards updating and adding new animations is relatively straightforward.

### Hue App Part 1

To control your Hue lights via an Alexa skill, you will need to register as a developer with Hue and create a Hue 'app'.

1. Head over to the [Hue developer login page](https://developers.meethue.com/login/) and create an account.
1. Once signed in, go to [your apps](https://developers.meethue.com/my-apps/) and click "Add new Remote Hue API app".
1. Fill out the required fields of the form:
   - `App name`: Any value you want
   - `Callback URL`: Fill this with `https://www.example.com`. We'll come back to this later.
   - `Application Description`: Anything you want
1. Click "Submit"

### Alexa Skill Part 1

1. Go to the Alexa Skill [development portal](https://developer.amazon.com/alexa/console/ask?). Create an account if you
   haven't yet, using the same email as your Amazon account.
1. Click "Create Skill"
1. Fill in a skill name and click the "Smart Home" tile. Click "Create Skill" in the top right corner.
1. Under step two "Smart Home service endpoint" of the next page, copy the field value of "Your Skill ID"
1. Keep this tab open, we'll come back to it.

### Set up SAM (AWS Serverless Application Model)

This project is built on top of [SAM](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/what-is-sam.html). This is due the many great benefits it offers, such as easy deployment and configuration, as well as
the capability to run your functions locally.

1. Follow the steps listed here, and [install and set up the SAM CLI](https://docs.aws.amazon.com/serverless-application-modellatest/developerguide/serverless-sam-cli-install.html).
   This will take you through the CLI installation, including creating an AWS account and creating an Admin user with
   proper IAM roles if you haven't already done so.
1. [Set up your AWS credentials](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-getting-started-set-up-credentials.html). You may leave the defaults as they are, we'll change that up later.

### Set up Git Repository

1. Clone this repository (or better yet fork then clone so you can contribute your own improvements or animations).
1. In [the template file](./template.yaml) replace `<your skill id here>` with the skill id obtained in Alexa Skill Part
   1 step 4.
1. Run `make build` in your terminal in the directory of the repository.
1. Run `sam deploy --guided`, saying yes to all prompts **and setting 'us-east-1' as the region (if you are in NA)**.
1. Once the deployment is complete, make a note of the values of `Alexa` and `AuthenticationAPI` in the output.

### Alexa Skill Part 2

1. Copy the `ClientId` and `ClientSecret` from your Hue app's card.
1. Navigate back to the Alexa Skill tab from Alexa Skill Part 1.
1. On the "Smart Home" tab, fill the "Default Endpoint" field with the ARN from the `Alexa` output obtained from the
   step above.
1. On the Account Linking page, we can fill out these fields using your Hue Apps information:

   1. Authorization URI: The URI to kick off the oAuth process with Hue. This is a URL with the following form:
   
      ```
      https://api.meethue.com/oauth2/auth?clientid=<your_hue_clientid>&response_type=code&state=<random string>&appid=<your hue app name>&deviceid=<any unique name>&devicename=<any unique name>

      ```

   1. `Access Token URI`: The URI to get access tokens from Hue. Use the `AuthenticationAPI` value obtained from the
      step above.
   1. `Your Client ID`: The `ClientId` from your Hue app.
   1. `Your Secret`: The `ClientSecret` from your Hue app.
   1. `Your Authentication Scheme`: HTTP Basic
   1. `Scope`: profile

1. Copy one of the Alexa Redirect URLs following depending on your region:
   - North America: `https://pitangui.amazon.com/…`
   - EU and India: `https://layla.amazon.com/...`
   - Far East: `https://alexa.amazon.co.jp...`

### Hue App Part 2

1. Change your Hue App's "Callback URL" to be the URL copied in the step above.

### Install the Skill

1. Open the Amazon Alexa app on your phone.
1. Navigate to Skills & Games > Your Skills
1. Under the "Dev" dropdown, you should see your newly created skill. Tap it and follow the account link instructions.
1. Discover devices. This should add all defined devices in the [alexa file](./lambdas/alexa/main.go) to your Alexa.
1. You're done, have fun customizing your Hue lights 🎉

## Future Plans

This is a rough first pass at this. Many things could be improved:

1. Make things in the `sunrise` device configurable from the Amazon Alexa app, such as duration, colors, etc.
1. Improve error communication to deliver errors to Alexa via speech output.
1. Refactor code to be able to host shared skill everyone can use, rather than having people deploy their own.
