import gflags
import httplib2
import os
import sys
import time
from time import sleep
from datetime import datetime
import dateutil.relativedelta
from rfc3339 import rfc3339

from apiclient.discovery import build
from oauth2client.file import Storage
from oauth2client.client import OAuth2WebServerFlow
from oauth2client.tools import run

FLAGS = gflags.FLAGS


# First, authenticate the application
print '========================================'
print '             vestige/1.0'
print '========================================'
print ''

# We have to check for three environment variables
# client_id, client_secret, and developer_key

print '** Checking for application keys...'

clientId = os.environ.get('VESTIGE_CLIENT_ID')
clientSecret = os.environ.get('VESTIGE_CLIENT_SECRET')
developerKey = os.environ.get('VESTIGE_DEVELOPER_KEY')

if clientId == None:
  print '!! You must set the VESTIGE_CLIENT_ID environment variable.'
if clientSecret == None:
  print '!! You must set the VESTIGE_CLIENT_SECRET environment variable.'
if developerKey == None:
  print '!! You must set the VESTIGE_DEVELOPER_KEY environment variable.'

if clientId == None or clientSecret == None or developerKey == None:
  print '!! You can find these developer keys on the Google API Console:'
  print '!! https://code.google.com/apis/console/'
  exit()


# Now, we can authenticate to Google

print '** Authenticating to Google...'

FLOW = OAuth2WebServerFlow(
    client_id = clientId,
    client_secret = clientSecret,
    scope = 'https://www.googleapis.com/auth/calendar',
    user_agent = 'vestige/1.0')

FLAGS.auth_local_webserver = False

storage = Storage('calendar.dat')
credentials = storage.get()
if credentials is None or credentials.invalid == True:
  credentials = run(FLOW, storage)

http = httplib2.Http()
http = credentials.authorize(http)

service = build(serviceName='calendar', version='v3', http=http, developerKey = developerKey)

print '** Authentication finished.'
print '** Ready.'


# workActive variable defines whether we're actively doing work
workActive = False

# Define convenience functions
def rewriteLine():
  sys.stdout.write('\r')
  sys.stdout.flush()

def timeLoop(startTime):
  while workActive == True:
    rewriteLine()
    delta = dateutil.relativedelta.relativedelta(startTime, datetime.now())
    print "Working for %d hours, %d minutes and %d seconds" % (rd.hours, rd.minutes, rd.seconds)


# Now we're authenticated. Define the loop.

def appLoop():
  print ''
  print '-- NEW WORK ITEM --------------------------'
  print ' * What are you working on?'
  workItem = raw_input('   ')

  # Now we have the work item, record the current time 
  startTime = datetime.now()
  workActive = True

  # Echo the start time
  print ''
  print ' * Started at ' + str(startTime.hour) + ':' + str(startTime.minute)

  # I tried to put in a "time counter" here, but it would have required ncurses and
  # threading and some other nasty stuff... will figure it out later

  # Wait for user to complete the task
  raw_input('   Hit Enter to finish work')

  # User has now finished task, calculate current time and enter into Google Calendar
  endTime = datetime.now()
  workActive = False

  # Construct event object
  event = {
      'summary': workItem,
      'start': {
        'dateTime': rfc3339(startTime)
      },
      'end': {
        'dateTime': rfc3339(endTime)
      }
    }
  
  print ''
  print ' * Sending event to Google...'
  
  # Send to Google Calendar
  eventSent = service.events().insert(calendarId = 'primary', body=event).execute()
  
  print ' * Event ' + eventSent['id'] + ' created. Done.'
  print '-- END WORK ITEM --------------------------'
  print ''

  # And let's do it again
  appLoop()

appLoop()
