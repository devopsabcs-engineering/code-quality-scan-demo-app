@description('Name prefix for storage resources')
param namePrefix string = 'cqscan'

@description('Azure region for all resources')
param location string = resourceGroup().location

var suffix = uniqueString(resourceGroup().id)
var storageAccountName = 'st${suffix}'
var containerName = 'scan-results'

resource storageAccount 'Microsoft.Storage/storageAccounts@2023-05-01' = {
  name: storageAccountName
  location: location
  kind: 'StorageV2'
  sku: {
    name: 'Standard_LRS'
  }
  properties: {
    isHnsEnabled: true
    accessTier: 'Hot'
    supportsHttpsTrafficOnly: true
    minimumTlsVersion: 'TLS1_2'
    allowBlobPublicAccess: false
    networkAcls: {
      defaultAction: 'Allow'
      bypass: 'AzureServices'
    }
  }
}

resource blobService 'Microsoft.Storage/storageAccounts/blobServices@2023-05-01' = {
  parent: storageAccount
  name: 'default'
}

resource container 'Microsoft.Storage/storageAccounts/blobServices/containers@2023-05-01' = {
  parent: blobService
  name: containerName
  properties: {
    publicAccess: 'None'
  }
}

output storageAccountName string = storageAccount.name
output containerName string = containerName
output storageAccountId string = storageAccount.id
output dfsEndpoint string = 'https://${storageAccount.name}.dfs.core.windows.net'
