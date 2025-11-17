/**
 * Validation utility i18n configuration - Normal theme
 */
export const validation = {
  errors: {
    required: 'This field is required',
    type: 'Field type is incorrect',
    min: 'Value too small, minimum is {min}',
    max: 'Value too large, maximum is {max}',
    length: 'Length does not meet requirements',
    pattern: 'Format is incorrect',
    email: 'Please enter a valid email address',
    phone: 'Please enter a valid phone number',
    url: 'Please enter a valid URL',
    username: 'Username can only contain letters, numbers, underscores and hyphens, length 3-20 characters',
    password: 'Password must contain uppercase and lowercase letters and numbers, at least 8 characters',
    custom: 'Validation failed',
    mustBeString: 'Must be string type',
    mustBeNumber: 'Must be number type',
    mustBePositive: 'Must be positive number',
    minLength: 'Minimum {n} characters',
    maxLength: 'Maximum {n} characters',
  },
  password: {
    strength: {
      weak: 'Weak',
      medium: 'Medium',
      strong: 'Strong',
      veryStrong: 'Very Strong',
    },
    feedback: {
      minLength: 'Password must be at least 8 characters',
      requireLowercase: 'Requires lowercase letters',
      requireUppercase: 'Requires uppercase letters',
      requireNumber: 'Requires numbers',
      suggestSpecial: 'Suggest including special characters',
      avoidRepeat: 'Avoid consecutive repeated characters',
    },
  },
  time: {
    days: '{n} days',
    hours: '{n} hours',
    minutes: '{n} minutes',
    seconds: '{n} seconds',
    daysHours: '{d} days {h} hours',
    hoursMinutes: '{h} hours {m} minutes',
    minutesSeconds: '{m} minutes {s} seconds',
  },
}
