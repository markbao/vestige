# vestige

### Command-line time tracking to Google Calendar, with task categorization. Written in Go.

Time passes by quickly, and sometimes we don't know what happened during our day. Vestige allows you to track what you do during your day and send that data to Google Calendar.

Partially inspired by Chris Dancy's insane [Google Calendar tracking system](http://www.wired.com/wiredenterprise/2013/02/quantified-work/all/).

## How?

First you...

````
❯ go run vestige.go --clientid="your-client-id.apps.googleusercontent.com" --secret="your-api-secret"

-- vestige 1.0 ----------------------------
 * Authenticating to Google...
 * Loading calendars...
 * Ready.


-- NEW WORK ITEM --------------------------
 * What are you working on?
   Relationships — Email John about new strategy

 * Started at 7:45PM
   Hit Enter to finish work

 * Sending to Google...
 * Sent.
-- END WORK ITEM --------------------------
````

And then...

![](http://i.imgur.com/k895ZlZ.png)

And then you get back to work.

## How it works

You enter `Notes - Wikipedia page on Emu War` and hit enter. Notes is the category, the rest is the activity.

````
-- NEW WORK ITEM --------------------------
 * What are you working on?
   Notes - Wikipedia page on Emu War

 * Started at 6:53AM
   Hit Enter to finish work
````

Once you're done, hit enter.

* The work item will be put into your Google Calendar.
* The duration will be your work time.
* It'll be categorized under the calendar *Notes*.
* (If that calendar doesn't exist, it'll create the calendar for you.)
* (If you just type `Wikipedia page on Emu War` without a category, it'll put it in your primary calendar.)

You can change the behavior of this, too.

* To change your default calendar, use `--default=NameOfCalendar`. All uncategorized notes will go in here.
* To restrict *all* of your work items to one calendar, use `--single=true`. This will put them in your default calendar.

## Install

First, install the Go packages for OAuth and the Google API Go Client:

````
go get code.google.com/p/goauth2/oauth
go get code.google.com/p/google-api-go-client/calendar/v3
````

Next, register an application at the
[Google API Console](https://code.google.com/apis/console/), and obtain a set
of OAuth consumer keys. Get your client ID and secret ready.

Execute the program like so:

````
go run vestige.go --clientid="your-client-id" --secret="your-api-secret"
````

Follow the on-screen instructions to link the app to your Google Calendar account.

#### Hey, what's this Python file doing here?

The first version of Vestige was written in Python. It's still usable, but the Go version is much more featureful.

## Todo

### Categorization

* Done! — Support different calendars (as categories). 'Personal - Doing taxes' would go into the Personal calendar.
* Done! — When a category doesn't exist, create it. (We'll need to load a list of calendars in the beginning, and load it again when we create one.)
* Done! — Allow a user to switch between creating categories when the hyphen is used, or just putting it into the primary calendar.

### Offline stuff

* Detect internet connection status.
* Save work items when offline, and submit when submitting a work item when online.

### OMG, data

* Add a script to tally up all the categories of any day, to see what 
* Give this script a graph (graphs rock)
* Allow this script to export HTML graphs that use JS charts (whoa man)
* Allow this script to create time-based graphs, maybe even stacked graphs, to make sense of how we're using our time during the day (serious stuff right here)

### Other cool stuff
* Make it possible to load the Client ID and Secret from the bash configuration or save it somewhere
* Add counter like 'Working for 1 hour and 25 minutes.'
* Figure out a way to make sending the event data to Google not so slow (is it doing more handshakes than necessary?)
* Make it pester you (ring a terminal bell?) when you don't update it for a few minutes.

