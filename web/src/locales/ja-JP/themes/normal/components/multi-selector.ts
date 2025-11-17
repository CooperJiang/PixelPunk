/**
 * Multi-selector component
 */
export const multiSelector = {
  addOption: 'オプション追加',
  deleteOption: 'オプション削除',
  confirmAdd: '追加確認',
  cancel: 'キャンセル',
  placeholders: {
    basic: '例：30m、2h、1d',
    withPermanent: '例：30m、2h、1d、permanent',
  },
  duration: {
    permanent: '永続',
    minute: '分',
    hour: '時間',
    day: '日',
  },
  hints: {
    format: 'サポート形式：30m（分）、1h（時間）、7d（日）、最小1分',
    permanent: ' · permanent（永続）',
  },
  errors: {
    empty: '期間を入力してください',
    guestNoPermanent: 'ゲストは永続ストレージをサポートしていません',
    invalidFormat: '無効な形式、サポート：30m（分）/1h（時間）/7d（日）',
    mustBePositive: '期間は0より大きい必要があります',
    minutesMin: '分は1未満にすることはできません',
    minutesMax: '分は最大365日',
    hoursMax: '時間は最大365日',
    guestHoursMax: 'ゲストは最大720時間（30日）',
    daysMax: '日は最大365日',
    guestDaysMax: 'ゲストは最大30日',
    exists: 'この期間は既に存在します',
  },
}

