/**
 * Multi-selector component
 */
export const multiSelector = {
  addOption: 'Add Option',
  deleteOption: 'Delete Option',
  confirmAdd: 'Confirm Add',
  cancel: 'Cancel',
  placeholders: {
    basic: 'e.g., 30m, 2h, 1d',
    withPermanent: 'e.g., 30m, 2h, 1d, permanent',
  },
  duration: {
    permanent: 'Permanent',
    minute: 'minutes',
    hour: 'hours',
    day: 'days',
  },
  hints: {
    format: 'Supported formats: 30m(minutes), 1h(hours), 7d(days), minimum 1 minute',
    permanent: ' · permanent(permanent)',
  },
  errors: {
    empty: 'Please enter duration',
    guestNoPermanent: 'Guests do not support permanent storage',
    invalidFormat: 'Invalid format, supports: 30m(minutes)/1h(hours)/7d(days)',
    mustBePositive: 'Duration must be greater than 0',
    minutesMin: 'Minutes cannot be less than 1',
    minutesMax: 'Minutes maximum 365 days',
    hoursMax: 'Hours maximum 365 days',
    guestHoursMax: 'Guests maximum 720 hours (30 days)',
    daysMax: 'Days maximum 365 days',
    guestDaysMax: 'Guests maximum 30 days',
    exists: 'This duration already exists',
  },
}
