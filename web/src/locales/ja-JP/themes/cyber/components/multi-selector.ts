/**
 * Multi Selector Component
 */
export const multiSelector = {
  addOption: 'オプションを注入',
  deleteOption: 'コマンドを削除',
  confirmAdd: '注入確認',
  cancel: '中止',
  placeholders: {
    basic: '例: 30m, 2h, 1d',
    withPermanent: '例: 30m, 2h, 1d, permanent',
  },
  duration: {
    permanent: '永続保存',
    minute: '分',
    hour: '時間',
    day: '日',
  },
  hints: {
    format: '30m(分) / 1h(時間) / 7d(日)をサポート、最小1分',
    permanent: ' · permanent(永続)',
  },
  errors: {
    empty: '期間パラメータを入力してください',
    guestNoPermanent: 'ゲストモードでは永続保存が無効です',
    invalidFormat: 'フォーマットエラー、30m / 1h / 7dである必要があります',
    mustBePositive: '期間は0より大きい必要があります',
    minutesMin: '分は1未満にできません',
    minutesMax: '分の最大制限は365日です',
    hoursMax: '時間の最大制限は365日です',
    guestHoursMax: 'ゲストの最大制限は720時間（30日）です',
    daysMax: '日の最大制限は365日です',
    guestDaysMax: 'ゲストの最大制限は30日です',
    exists: 'この期間は既に存在します',
  },
}
