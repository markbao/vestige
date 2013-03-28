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

## Todo

### Categorization

* Support different calendars (as categories). 'Personal - Doing taxes' would go into the Personal calendar.
* When a category doesn't exist, create it. (We'll need to load a list of calendars in the beginning, and load it again when we create one.)
* Allow a user to switch between creating categories when the hyphen is used, or just putting it into the primary calendar.
* Allow a user to put 

### Offline stuff

* Detect internet connection status.
* Save work items when offline, and submit when submitting a work item when online.

### OMG, data

* Add a script to tally up all the categories of any day, to see what 
* Give this script a graph (graphs rock)
* Allow this script to export HTML graphs that use JS charts (whoa man)
* Allow this script to create time-based graphs, maybe even stacked graphs, to make sense of how we're using our time during the day (serious stuff right here)

### Other cool stuff
* Add counter like 'Working for 1 hour and 25 minutes.'
* Figure out a way to make sending the event data to Google not so slow (is it doing more handshakes than necessary?)
* Fix bug where single minutes show up as 2:5 (not 2:05).
* Make it pester you (ring a terminal bell?) when you don't update it for a few minutes. (Requires another thread since python is waiting at raw_input on the "What are you working on?" prompt.)
