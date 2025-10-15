# gospectrum

Test by Storwize V7000 (software version 8.3.x)

### ArrayInstance

#### MdiskId

The identity of the array MDisk.  
> field: mdisk_id  
> type: string

#### MdiskName

The name of array Mdisk.  
> field: mdisk_name  
> type: string

#### Status

 Indicates the array status.
> field: status  
> type: string<enum(float64)>
> - offline: 0
> - online: 1
> - degraded: 2
> - degraded_paths: 3
> - degraded_ports: 4
