# Tumult
Implementation of a Tumult connector, using library [Tumult Analytics](https://gitlab.com/tumult-labs/analytics).

## Tumult Supported functions
| Measurements   | Implemented | Supported by WebDP  | 
| :------------: | :---------: | :-----------------: |  
| Count          | Implemented | Supported           | 
| Sum            | Implemented | Supported           |
| Mean           | Implemented | Supported           |
| Min            | Implemented | Supported           |
| Max            | Implemented | Supported           |

| Transformations  | Implemented | Supported by WebDP  |  
| :--------------: | :---------: | :-----------------: | 
| Rename           | Implemented | Supported           |
| Filter           | Implemented | Supported           |
| Select           | Implemented | Supported           |
| Map              | Implemented | Supported           |
| Bin              | Implemented | Supported           |
| GroupBy          | Implemented | Supported           |

| Mechanisms             | Implemented | Supported by WebDP  | 
| :--------------------: | :---------: | :-----------------: | 
| Laplace mechanism      | Implemented | Supported           | 
| Gaussian mechanism     | Implemented | Supported           |

### Information
For an aggregate function (a query step that ends with "Measurement") to be valid, the user has to specify which 
column in the dataset to aggregate as well as which noise mechanism to use. If the PrivacyNotion is PureDP, then 
the Laplace mechanism is supported. If the PrivacyNotion is ApproxDP, then the Gauss mechanism is supported. Although the 
choice of mechanism is at this point in time forced, it can come to change in the future, which is why we prompt the user 
to be explicit. 

## Endpoints
| Endpoints      | Implemented | 
| :------------: | :---------: | 
| /evaluate      | Implemented | 
| /validate      | Implemented | 
| /accuracy      | False       | 
| /documentation | Implemented | 
| /functions     | Implemented |
| /cache/:id     | Implemented |

In the current version of Tumult library there are no available accuracy functions, hence it has not been implemented in this connector. 

## Creator
Created by David al Amiri with inspiration from Dpella's WebDP v1
