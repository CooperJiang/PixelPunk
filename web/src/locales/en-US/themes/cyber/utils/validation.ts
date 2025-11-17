/**
 * Validation Utils - Cyber Theme
 */
export const validation = {
  errors: {
    required: 'This field is required',
    type: 'Field type incorrect',
    min: 'Value too small, minimum is {min}',
    max: 'Value too large, maximum is {max}',
    length: 'Length does not meet specification',
    pattern: 'Format validation failed',
    email: 'Enter valid email address',
    phone: 'Enter valid phone number',
    url: 'Enter valid URL address',
    username: 'Username can only contain letters, numbers, underscores and hyphens, length 3-20',
    password: 'Password must contain uppercase, lowercase and numbers, at least 8 characters',
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
      veryStrong: 'Very strong',
    },
    feedback: {
      minLength: 'Password length at least 8 characters',
      requireLowercase: 'Requires lowercase letters',
      requireUppercase: 'Requires uppercase letters',
      requireNumber: 'Requires numbers',
      suggestSpecial: 'Recommend including special characters',
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
