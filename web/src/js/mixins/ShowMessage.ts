import Snackbar from 'node-snackbar'
import { Component, Vue } from 'vue-property-decorator'

@Component
export default class ShowMessage extends Vue {
  protected showMessage (message: string): void {
    Snackbar.show({
      duration: 3000,
      showAction: false,
      text: message
    })
  }
}