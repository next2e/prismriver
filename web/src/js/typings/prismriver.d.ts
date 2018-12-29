interface IMedia {
  ID: string
  CreatedAt: Date
  UpdatedAt: Date
  DeletedAt: Date
  Length: number
  Title: string
  Type: string
}

interface IQueueItem {
  Downloading: boolean
  DownloadProgress: number
  Media: IMedia
}
