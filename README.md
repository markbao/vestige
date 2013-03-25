# vestige

### Command-line time tracking to Google Calendar.

Time passes by quickly, and sometimes we don't know what happened during our day. Vestige allows you to track what you do during your day and send that data to Google Calendar.

Partially inspired by Chris Dancy's insane [Google Calendar tracking system](http://www.wired.com/wiredenterprise/2013/02/quantified-work/all/).

## How?

First you...

````
❯ python vestige.py
========================================
             vestige/1.0
========================================

** Checking for application keys...
** Authenticating to Google...
** Authentication finished.
** Ready.

-- NEW WORK ITEM --------------------------
 * What are you working on?
   Vestige - Writing README Markdown doc

 * Started at 16:7
   Hit Enter to finish work

 * Sending event to Google...
 * Event if0e8i8pi1eee6f4at1oals7uc created. Done.
-- END WORK ITEM --------------------------
````

And then...

![](http://i.imgur.com/PanPwos.png)

And then you get back to work.

## Install

You need to install the python modules `google-api-python-client` and `rfc3339`.

Then you download this package and `python vestige.py`.

It'll tell you to set the environment variables with your Google API keys.

To do this, go to the [Google API Console](https://code.google.com/apis/console/)
and register an application. Register it for Google Calendar use. Then, set it up
for OAuth for a downloadable application. Grab your application keys and OAuth
keys and set the environment variables for `VESTIGE_CLIENT_ID`, 
`VESTIGE_CLIENT_SECRET`, and `VESTIGE_DEVELOPER_KEY`.

Then do `python vestige.py` again and follow the on-screen instructions.

Now you're done.
