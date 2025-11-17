/**
 * Multi Selector Component
 */
export const multiSelector = {
  addOption: 'Inject Option',
  deleteOption: 'Remove Command',
  confirmAdd: 'Confirm Inject',
  cancel: 'Abort',
  placeholders: {
    basic: 'Example: 30m, 2h, 1d',
    withPermanent: 'Example: 30m, 2h, 1d, permanent',
  },
  duration: {
    permanent: 'Permanent Storage',
    minute: 'minute',
    hour: 'hour',
    day: 'day',
  },
  hints: {
    format: 'Supports 30m(minutes) / 1h(hours) / 7d(days), minimum 1 minute',
    permanent: ' · permanent(permanent)',
  },
  errors: {
    empty: 'Please enter duration parameter',
    guestNoPermanent: 'Guest mode disables permanent storage',
    invalidFormat: 'Format error, should be 30m / 1h / 7d',
    mustBePositive: 'Duration must be greater than 0',
    minutesMin: 'Minutes cannot be less than 1',
    minutesMax: 'Minutes max limit is 365 days',
    hoursMax: 'Hours max limit is 365 days',
    guestHoursMax: 'Guest max limit is 720 hours (30 days)',
    daysMax: 'Days max limit is 365 days',
    guestDaysMax: 'Guest max limit is 30 days',
    exists: 'This duration already exists',
  },
}
