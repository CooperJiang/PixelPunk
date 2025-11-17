/* 导入所有组件 - 统一使用默认导入格式 */
import Switch from './Switch'
import Dropdown from './Dropdown'
import Navbar from './Navbar/index.vue'
import Toast, { useToast } from './Toast'
import Button from './Button'
import Dialog from './Dialog'
import Input from './Input'
import Popconfirm from './Popconfirm'
import Breadcrumb from './Breadcrumb'
import { Loading } from './Loading'
import FolderTree from './FolderTree'
/* import cyberFileViewer from './cyberFileViewer' // deprecated: use ImageViewer */
import Pagination from './Pagination'
import UserAvatar from './UserAvatar'
import Logo from './Logo'
import { Radio, RadioGroup } from './Radio'
import { Slider } from './Slider'
import File from './File'
import FileDetailModal from './FileDetailModal'
import Tooltip from './Tooltip'
import DatePicker from './DatePicker'
import Checkbox from './Checkbox'
import MobileNavigation from './MobileNavigation'
import ShareFolder from './ShareFolder'
import { ShareFile } from './ShareFile'
import WaterfallLayout from './WaterfallLayout'
import FileViewer from './FileViewer'
import ParticleBackground from './ParticleBackground/index.vue'
import _Footer from './Footer/index.vue'
import CyberBackground from './Background'
import Drawer from './Drawer'
import CyberDataTable from './Table'
import CyberFileLoading from './FileLoading/index.vue'
import CyberResumableUploads from './ResumableUploads'
import Badge from './Badge/index.vue'
import CyberConfirmDialog from './ConfirmDialog'
import CyberIconButton from './IconButton'
import CyberIconPicker from './IconPicker'
import Copyright from './Copyright/index.vue'
import CommunityDialog from './CommunityDialog/index.vue'
import NotificationDialog from './NotificationDialog/index.vue'
import DuplicateBadge from './DuplicateBadge/index.vue'
import GlobalUploadFloat from './GlobalUploadFloat/index.vue'
import GlobalUploadDrawer from './GlobalUploadDrawer/index.vue'
import MultiSelector from './MultiSelector'
import WatermarkConfig from './WatermarkConfig/index.vue'
import FileExpiryTag from './FileExpiryTag'
import CyberAdminWrapper from './AdminWrapper'
import CyberSideNavTabs from './SideNavTabs'
import CyberSidebarNav from './CyberSidebarNav'
import CyberAccordion from './Accordion'
import CyberSkeleton from './CyberSkeleton'
import StatsCard from './StatsCard'
import StatsSection from './StatsSection'
import AccessLevelToggle from './AccessLevelToggle'
import Tag from './Tag'
import SmartTagContainer from './SmartTagContainer'
import CyberSidebar from './Sidebar/index.vue'
import LayoutSwitcher from './LayoutSwitcher/index.vue'
import LayoutToggleButton from './LayoutToggleButton/index.vue'
import PageContainer from './LayoutAdaptive/PageContainer.vue'
import GridContainer from './LayoutAdaptive/GridContainer.vue'
import CardContainer from './LayoutAdaptive/CardContainer.vue'
import CyberContextMenu from './CyberContextMenu'
import CyberHomeBackground from './CyberHomeBackground/index.vue'
import CyberThemeToggle from './ThemeToggle'
import Card from './Card'
import ColorPicker from './ColorPicker/index.vue'
import AnnouncementButton from './AnnouncementButton/index.vue'
import CyborgCharacter from './CyborgCharacter/index.vue'

/* 使用cyber前缀导出组件 */
export const components = {
  cyberCard: Card,
  cyberButton: Button,
  cyberToast: Toast,
  cyberDialog: Dialog,
  cyberSwitch: Switch,
  cyberPopconfirm: Popconfirm,
  cyberPagination: Pagination,
  cyberBreadcrumb: Breadcrumb,
  cyberDropdown: Dropdown,
  cyberLoading: Loading,
  cyberInput: Input,
  cyberFileViewer: FileViewer,
  cyberNavbar: Navbar,
  cyberFolderTree: FolderTree,
  cyberUserAvatar: UserAvatar,
  cyberRadio: Radio,
  cyberRadioGroup: RadioGroup,
  cyberLogo: Logo,
  cyberSlider: Slider,
  cyberFile: File,
  cyberFileDetailModal: FileDetailModal,
  cyberTooltip: Tooltip,
  cyberDatePicker: DatePicker,
  cyberCheckbox: Checkbox,
  cyberMobileNavigation: MobileNavigation,
  cyberShareFolder: ShareFolder,
  cyberShareFile: ShareFile,
  cyberWaterfallLayout: WaterfallLayout,
  cyberFileLoading: CyberFileLoading,
  cyberParticleBackground: ParticleBackground,
  cyberBackground: CyberBackground,
  cyberDrawer: Drawer,
  cyberDataTable: CyberDataTable,
  cyberResumableUploads: CyberResumableUploads,
  cyberBadge: Badge,
  cyberConfirmDialog: CyberConfirmDialog,
  cyberIconButton: CyberIconButton,
  cyberIconPicker: CyberIconPicker,
  cyberCopyright: Copyright,
  cyberCommunityDialog: CommunityDialog,
  cyberNotificationDialog: NotificationDialog,
  cyberDuplicateBadge: DuplicateBadge,
  cyberGlobalUploadFloat: GlobalUploadFloat,
  cyberGlobalUploadDrawer: GlobalUploadDrawer,
  cyberMultiSelector: MultiSelector,
  CyberMultiSelector: MultiSelector,
  cyberWatermarkConfig: WatermarkConfig,
  cyberFileExpiryTag: FileExpiryTag,
  cyberAdminWrapper: CyberAdminWrapper,
  cyberSideNavTabs: CyberSideNavTabs,
  cyberSidebarNav: CyberSidebarNav,
  cyberAccordion: CyberAccordion,
  cyberSkeleton: CyberSkeleton,
  cyberStatsCard: StatsCard,
  cyberStatsSection: StatsSection,
  cyberAccessLevelToggle: AccessLevelToggle,
  cyberTag: Tag,
  cyberSmartTagContainer: SmartTagContainer,
  cyberSidebar: CyberSidebar,
  cyberLayoutSwitcher: LayoutSwitcher,
  cyberLayoutToggleButton: LayoutToggleButton,
  cyberPageContainer: PageContainer,
  cyberGridContainer: GridContainer,
  cyberCardContainer: CardContainer,
  cyberContextMenu: CyberContextMenu,
  cyberHomeBackground: CyberHomeBackground,
  cyberThemeToggle: CyberThemeToggle,
  cyberColorPicker: ColorPicker,
  cyberAnnouncementButton: AnnouncementButton,
  cyberCyborgCharacter: CyborgCharacter,
}

/* 导出useToast函数供全局使用 */
export { useToast }

/* 统一导出 CyberBadge */
export { default as CyberBadge } from './Badge/index.vue'
export { default as DuplicateBadge } from './Badge/index.vue'
export { default as SimilarityBadge } from './Badge/index.vue'

/* 组件映射表 */
const componentMap = {
  CyberCard: Card,
  CyberSwitch: Switch,
  CyberDropdown: Dropdown,
  CyberNavbar: Navbar,
  CyberToast: Toast,
  CyberButton: Button,
  CyberDialog: Dialog,
  CyberInput: Input,
  CyberPopconfirm: Popconfirm,
  CyberBreadcrumb: Breadcrumb,
  CyberLoading: Loading,
  CyberFolderTree: FolderTree,
  CyberFileViewer: FileViewer,
  CyberPagination: Pagination,
  CyberUserAvatar: UserAvatar,
  CyberRadio: Radio,
  CyberRadioGroup: RadioGroup,
  CyberLogo: Logo,
  CyberSlider: Slider,
  CyberFile: File,
  CyberFileDetailModal: FileDetailModal,
  CyberTooltip: Tooltip,
  CyberDatePicker: DatePicker,
  CyberCheckbox: Checkbox,
  CyberMobileNavigation: MobileNavigation,
  CyberShareFolder: ShareFolder,
  CyberShareFile: ShareFile,
  CyberWaterfallLayout: WaterfallLayout,
  CyberFileLoading,
  CyberParticleBackground: ParticleBackground,
  CyberBackground,
  CyberDrawer: Drawer,
  CyberDataTable,
  CyberTable: CyberDataTable,
  CyberResumableUploads,
  CyberBadge: Badge,
  CyberConfirmDialog,
  CyberIconButton,
  CyberIconPicker,
  CyberCopyright: Copyright,
  CyberCommunityDialog: CommunityDialog,
  CyberNotificationDialog: NotificationDialog,
  CyberDuplicateBadge: DuplicateBadge,
  CyberGlobalUploadFloat: GlobalUploadFloat,
  CyberGlobalUploadDrawer: GlobalUploadDrawer,
  CyberMultiSelector: MultiSelector,
  CyberWatermarkConfig: WatermarkConfig,
  CyberFileExpiryTag: FileExpiryTag,
  CyberAdminWrapper,
  CyberSideNavTabs,
  CyberSidebarNav,
  CyberAccordion,
  CyberSkeleton,
  CyberStatsCard: StatsCard,
  CyberStatsSection: StatsSection,
  CyberAccessLevelToggle: AccessLevelToggle,
  CyberTag: Tag,
  CyberSmartTagContainer: SmartTagContainer,
  CyberSidebar,
  CyberLayoutSwitcher: LayoutSwitcher,
  CyberLayoutToggleButton: LayoutToggleButton,
  CyberPageContainer: PageContainer,
  CyberGridContainer: GridContainer,
  CyberCardContainer: CardContainer,
  CyberContextMenu,
  CyberHomeBackground,
  CyberThemeToggle,
  CyberColorPicker: ColorPicker,
  CyberAnnouncementButton: AnnouncementButton,
  CyborgCharacter,
  CyberCyborgCharacter: CyborgCharacter,
}

/* 组件库插件安装函数 */
import type { App, Plugin } from 'vue'

/* 支持按需注册的组件库 */
export const createCyberUI = (options: { components?: string[] } = {}): Plugin => ({
  install(app: App) {
    const { components = [] } = options

    if (components.length > 0) {
      components.forEach((name) => {
        if (componentMap[name as keyof typeof componentMap]) {
          app.component(name, componentMap[name as keyof typeof componentMap])
        }
      })
    } else {
      Object.entries(componentMap).forEach(([name, component]) => {
        app.component(name, component)
      })
    }
  },
})

const CyberComponentsPlugin: Plugin = createCyberUI()

export default CyberComponentsPlugin
