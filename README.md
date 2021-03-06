# Vestige

## End time ignorance.

**Time passes by quickly, and sometimes a day passes by and we don't know how it happened. We work on a task, and then another, and then the day is almost over—and we don't know that we spent 30 minutes replying to that email, or 45 minutes reading the news.**

Time slips by, and we're oblivious to its passing because we don't quantify it. We need to quantify the things we do during the day and the time it takes us to do them, so we can study it and make changes to our behavior, to manage our time better.

The goal is to use software to track what tasks and projects you work on during the day, and then look back on your day to see what exactly took up your time.

Now, you go from *oblivious* to what took up time during the day, to now being *aware* of what is taking up your time, and from there, you can *decide* to make changes as necessary.

## What is it?

Vestige tracks your activities and puts them in Google Calendar. You start the app, type in what you're working on, and when you're done, hit Enter to finish the work item. That time block will be sent to Google Calendar so you can look at it later; now, you can start a new work item.

It supports task categorization into calendars, as well as a terminal bell to notify you when you've been idle. Written in Go.

Partially inspired by Chris Dancy's insane [Google Calendar tracking system](http://www.wired.com/wiredenterprise/2013/02/quantified-work/all/).

![](http://i.imgur.com/k895ZlZ.png)

## How it works

Run the program.

````
go run vestige.go --clientid="your-client-id" --secret="your-api-secret"
````

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

* **Default calendar**: To change your default calendar, use `--default=NameOfCalendar`. All uncategorized notes will go in here.
* **Single calendar mode**: To restrict *all* of your work items to one calendar, use `--single`. This will put them in your default calendar.
* **Idle reminder**: To be notified when you've been idle for more than 2 minutes, use `--remind`. This will sound a terminal bell when you've been idle for 2 minutes.

## Install

First, install the Go packages for OAuth and the Google API Go Client:

````
go get code.google.com/p/goauth2/oauth
go get code.google.com/p/google-api-go-client/calendar/v3
````

You may need to use `sudo` to get these installed into your Go `pkg` directory.

Next, register an application at the
[Google API Console](https://code.google.com/apis/console/), enable Google
Calendar on your app, and obtain a set of OAuth consumer keys. Get your
client ID and secret ready.

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
* Allow tab-completion of known categories
* Integrate spelling suggestions for misspelled categories
* Trim calendar names (don't allow preceding or following spaces)

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
* Done! — Make it pester you (ring a terminal bell?) when you don't update it for a few minutes.
* Done! — Add option to cancel current work item

