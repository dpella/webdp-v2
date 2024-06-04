# GoogleDP
Implementation of a GoogleDP connector, using library [GoogleDP](https://github.com/google/differential-privacy).

## GoogleDP Supported functions
| Measurements       | Implemented | Supported by WebDP  | 
| :----------------: | :---------: | :-----------------: |  
| Count              | Implemented | Supported           | 
| Sum                | Implemented | Supported           |
| Mean               | Implemented | Supported           |
| Variance           | Implemented | Not Supported       | 
| Standard deviation | Implemented | Not Supported       | 
| Quantiles          | Implemented | Not Supported       | 

| Transformations  | Implemented | Supported by WebDP  |  
| :--------------: | :---------: | :-----------------: | 
| Filter           | Implemented | Supported           |
| Bin              | Implemented | Supported           |

| Mechanisms             | Implemented | Supported by WebDP  | 
| :--------------------: | :---------: | :-----------------: | 
| Laplace mechanism      | Implemented | Supported           | 
| Gaussian mechanism     | Implemented | Supported           |

## Endpoints
| Endpoints      | Implemented | 
| :------------: | :---------: | 
| /evaluate      | Implemented | 
| /validate      | Implemented | 
| /accuracy      | False       | 
| /documentation | Implemented | 
| /functions     | Implemented |
| /cache/:id     | False       |

In the current version of GoogleDP library there are no available accuracy functions, hence it has not been implemented in this connector. 

## Added Functions
In this connector implementation filtering and binning has been implemented as a complement to the GoogleDP library functions.

### Filtering
The available filtering operations are: "<", ">", "<=", ">=", "==", "!="
  
If filtering numbers all operations are available.
If filtering text only operators "==", "!=" are available.

Conditions for filtering:
Each filter needs to be seperately written. Example ["ex_col > 20", "ex_col < 50"]
Filters are strictly AND. Example ["ex_col > 20", "ex_col < 50"] will return the rows that are 20 < x < 50.
Example: ["ex_col < 20", "ex_col > 50"] will return empty.

### Bins
Binning has been implemented but in its current configuration only works for integers

The correct format for a binning operation is an integer array with unique values and in ascending order.
Example: [10, 20, 30, 40, 50] or [10, 50, 80]
Example of wrong format: [50, 20, 60, 50]

## Limitations
Due to the limited support of transformation functions in the GoogleDP library there are a very limited support for transformations in the current connector. Filtering is limited and binning is even more limited compared to the Tumult connector.

## Additional work
* Increase test coverage
* Implement more transformations that are DP secure, such as increasing scope of binning and filtering while ensuring DP.

## Creator
Created by David al Amiri
