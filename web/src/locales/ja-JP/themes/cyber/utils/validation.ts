/**
 * Validation Utils - Cyber Theme
 */
export const validation = {
  errors: {
    required: 'このフィールドは必須です',
    type: 'フィールドタイプが正しくありません',
    min: '値が小さすぎます。最小値は{min}です',
    max: '値が大きすぎます。最大値は{max}です',
    length: '長さが仕様を満たしていません',
    pattern: 'フォーマット検証失敗',
    email: '有効なメールアドレスを入力してください',
    phone: '有効な電話番号を入力してください',
    url: '有効なURLアドレスを入力してください',
    username: 'ユーザー名は文字、数字、アンダースコア、ハイフンのみ使用可能、長さ3-20文字',
    password: 'パスワードは大文字、小文字、数字を含む必要があります。少なくとも8文字',
    custom: '検証失敗',
    mustBeString: '文字列タイプである必要があります',
    mustBeNumber: '数値タイプである必要があります',
    mustBePositive: '正の数である必要があります',
    minLength: '最小{n}文字',
    maxLength: '最大{n}文字',
  },
  password: {
    strength: {
      weak: '弱い',
      medium: '中程度',
      strong: '強い',
      veryStrong: '非常に強い',
    },
    feedback: {
      minLength: 'パスワード長は少なくとも8文字',
      requireLowercase: '小文字が必要です',
      requireUppercase: '大文字が必要です',
      requireNumber: '数字が必要です',
      suggestSpecial: '特殊文字を含めることを推奨',
      avoidRepeat: '連続した繰り返し文字を避けてください',
    },
  },
  time: {
    days: '{n}日',
    hours: '{n}時間',
    minutes: '{n}分',
    seconds: '{n}秒',
    daysHours: '{d}日{h}時間',
    hoursMinutes: '{h}時間{m}分',
    minutesSeconds: '{m}分{s}秒',
  },
}
