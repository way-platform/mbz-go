# Remote Maintenance Support API Reference

This is an easy step-by-step guide that shows you how to use the Remote Maintenance Support API. For advanced usage, you can find detailed information about parameters and their values in the [Specification](https://developer.mercedes-benz.com/apis/remote_maintenance_api/specification) area. For generating clients you can also download the openAPI specification (swagger files) there.

API Endpoint

```http
https://api.mercedes-benz.com/remotemaintenance/v3/
```

## Authentication

This API is secured with OAuth 2.0. You will need a valid client ID and client secret to use this API and obtain a valid access token for your end-users, which you need to include in your API calls. You can find detailed documentation on how to use APIs secured with OAuth 2.0 and obtain an access token on the [OAuth documentation](https://developer.mercedes-benz.com/content-page/oauth-documentation) page.

You can manage and view your client credentials in the [console](https://developer.mercedes-benz.com/console). Your client credentials carry many privileges, so be sure to keep them secure! **Do not share your secret client credentials** in publicly accessible areas such GitHub, client-side code, and so forth.

All API requests must be made over [HTTPS](https://en.wikipedia.org/wiki/HTTPS). Calls made over plain HTTP will fail. API requests without authentication will also fail.

Example Header

```bash
-H "authorization: Bearer <insert_the_access_token_here>"
```

### Required Scopes

For the OAuth calls, you need to incorporate the following scope when requesting authorization from the user:

- _mb:vehicle:rms:reader_ - Maintenance relevant data of your vehicle: The third party provider receives remote access to maintenance-relevant of your vehicles data (e.g. the actual mileage) over the air via the built in SIM card.


The [OAuth documentation](https://developer.mercedes-benz.com/content-page/oauth-documentation) page contains more details on how to request these scopes.

|   |   |
|---|---|
|Note|Remote Maintenance Support data retrieved by using this API is personal data within the meaning of the GDPR. Therefore, please adhere to the terms of use and data protection regulations within the general terms and conditions of the Remote Maintenance Support API when using this data.|

## Remote Maintenance Support

The Remote Maintenance Support API will provide the possibility for accredited 3rd party applications to access maintenance relevant vehicle data remotely on behalf of the Mercedes-Benz customer. To use the endpoints you need a valid VIN/FIN (**vehicleId**).

### 1. Get the Available Resources That Can Be Read Out

You can use the **resourceReadouts** resource of the Remote Maintenance Support API to get the resources that are available for readout. The response contains the available resources for the corresponding car. You can specify the car using the vehicleId path parameter in the request and verify the timestamp of the existing maintenance relevant data without using credits

**Content Description**
**id**: Unique Id of the readout
**asyncStatus**: Reflects the status of the readout
**receivedTimestamp**: Timestamp at which the data received the server. If the value is “-“ it means that there is no data available. You can than call the Vehicle Data Readout Refresh to get new data from the vehicle.
**messageTimestamp**: Timestamp of the API request
**vehicleId**: The vehicle identifier of the extended vehicle, e.g VIN
**resources**:
name: Reflects the name of the resource
version: Reflects the version of the resource
api: Reflects the URL of the /resourceReadouts endpoint

Example Request

```bash
curl
  -X POST "https://api.mercedes-benz.com/remotemaintenance/v3/vehicles/<insert_your_vehicle_id_here>/resourceReadouts"
  -H "accept: */*"
  -H "authorization: Bearer <insert_the_access_token_here>"
  -H "Content-Length:0"
```

Example Response

```json
{
  "resourceReadout": {
    "id": "exvs2510264021768535",
    "asyncStatus": "Complete",
    "receivedTimestamp": "2020-09-01T15:28:59+00:00",
    "messageTimestamp": "2020-09-08T07:26:24+00:00",
    "vehicleId": "WDD111111CAR01000",
    "resources": [
      {
        "name": "Value Readout",
        "version": "1.0.0.0",
        "api": "https://api.mercedes-benz.com/remotemaintenance/v3/vehicles/WDD111111CAR01000/valueReadouts"
      },
      {
        "name": "Warning Indicator Readout",
        "version": "1.0.0.0",
        "api": "https://api.mercedes-benz.com/remotemaintenance/v3/vehicles/WDD111111CAR01000/language/{language}/warningIndicatorReadouts"
      },
      {
        "name": "Vehicle Data Readout Refresh",
        "version": "1.0.0.0",
        "api": "https://api.mercedes-benz.com/remotemaintenance/v3/vehicles/WDD111111CAR01000/vehicleDataReadoutRefresh"
      }
    ]
  }
}
```

### 2. Read RMS Value Readout

You can use the **valueReadouts** resource: Read mileage, service code and workshop code for Mercedes-Benz vehicles.

**Content Description**
**asyncStatus**: Reflects the status of the readout
**id**: Unique Id of the readout
**messageTimestamp**: Timestamp of the API request
**receivedTimestamp**: Timestamp at which the data received the server
**vehicleId**: The vehicle identifier of the extended vehicle, e.g VIN
**odometerValue**: Mileage of the vehicle
**serviceCode**: The code for Cars and Vans to specify which maintenance type has to be done, e.g. A or B.
**workshopCode**: The code to define which maintenance steps have to be done.
**timeToService**: The remaining time in [days] to the next service.
**distanceToService**: The remaining distance in [km] to the next service. Only available for passenger cars and vans.

Example Request

```bash
curl
  -X POST "https://api.mercedes-benz.com/remotemaintenance/v3/vehicles/<insert_your_vehicle_id_here>/valueReadouts"
  -H "accept: */*"
  -H "authorization: Bearer <insert_the_access_token_here>"
  -H "Content-Length:0"
```

Example Response

```json
{
  "valuereadout": {
    "id": "exvs5719578132647805",
    "asyncStatus": "true",
    "messageTimestamp": "2020-07-01T14:14:07+00:00",
    "receivedTimestamp": "2020-07-01T14:14:07+00:00",
    "vehicleId": "WDD111111CAR01000",
    "maintenanceData": [
      {
        "odometerValue": "12153.0 km",
        "serviceCode": "A",
        "workshopCode": "505",
        "timeToService": "101.0 d",
        "distanceToService": "13091.0 km",
        "maintenanceData": [
          {
            "type": "Reifendruecke",
            "name": "Tirepressures",
            "line": [
              {
                "type": "TirePressureFrontLeft",
                "value": {
                  "doublevalue": 2.7,
                  "unit": "bar"
                }
              },
              { ... }
            ]
          }
        ]
      }
    ]
  }
}
```

### 3. Read RMS Warning Indicator Readout

You can use the **warningIndicatorReadouts** resource to read maintenance relevant warning lamps and event messages for passenger cars and vans.

**Content Description**
**asyncStatus**: Reflects the status of the readout
**id**: Unique Id of the readout
**messageTimestamp**: Timestamp of the API request
**receivedTimestamp**: Timestamp at which the data received the server
**vehicleId**: The vehicle identifier of the extended vehicle, e.g. VIN
**warningIndicatorReadouts_status**: Reflects the status of the Warning Indicator
**warningIndicatorReadouts_originalOdometerValue**: Mileage when the Warning Indicator occurred the first time
**warningIndicatorReadouts_mostRecentOdometerValue**: Mileage when the Warning Indicator occurred the last time
**warningIndicatorReadouts_frequencyCounter**: How many times the Warning Indicator occurred. If it’s 1.0, than originalOdometer and mostRecentOdometer are the same
**warningIndicatorReadouts_ignitionCycleCounter**: Indicates the number of consecutive operation cycles recorded since a Warning Indicator first became stored
**warningIndicatorReadouts_text_symptomText**: Headline of the Warning Indicator
**warningIndicatorReadouts_text_descriptionText**: The description of the Warning indicator of the digital manual with treatment recommendation
**warningIndicatorReadouts_icon**: If available, the icon as base64 encoded. Decode to get the icon

Example Request

```bash
curl
  -X POST "https://api.mercedes-benz.com/remotemaintenance/v3/vehicles/<insert_your_vehicle_id_here>/language/<insert_your_language_id>/warningIndicatorReadouts"
  -H "accept: */*"
  -H "authorization: Bearer <insert_the_access_token_here>"
  -H "Content-Length:0"
```

Example Response

```json
{
  "warningIndicatorReadout": {
    "id": "exvs4934615106146365",
    "asyncStatus": "Complete",
    "receivedTimestamp": "2020-09-01T15:28:59+00:00",
    "messageTimestamp": "2020-09-08T08:13:20+00:00",
    "vehicleId": "WDD111111CAR01000",
    "warningIndicators": [
      {
        "status": "STATUS_NOT_ACTIVE",
        "originalOdometerValue": "9904.0 km",
        "mostRecentOdometerValue": "9904.0 km",
        "frequencyCounter": "1",
        "ignitionCycleCounter": "3",
        "text": {
          "symptomText": "There is insufficient brake fluid in the brake fluid reservoir.",
          "descriptionText": "Example text"
        },
        "icon": "data:image/png;base64,example"
      }
    ]
}
```

### 4. Vehicle Data Readout Refresh

Optional function, which can be used in case of an outdated timestamp (“resourceReadouts” output), to receive a new vehicle data readout. The prerequisite for using the function for passenger cars is the additional activation of “Remote Diagnosis” in the Mercedes me Portal by the vehicle owner. The function requests a new data record from the vehicle data memory of passenger cars.

**Content Description**
**id**: Unique Id of the readout
**asyncStatus**: Reflects the status of the readout
**messageTimestamp**: Timestamp of the API request
**vehicleId**: The vehicle identifier of the extended vehicle, e.g VIN
**message**: A textual information that the request was successful or not

Example Request

```bash
curl
  -X POST "https://api.mercedes-benz.com/remotemaintenance/v3/vehicles/<insert_your_vehicle_id_here>/vehicleDataReadoutRefresh"
  -H "accept: */*"
  -H "authorization: Bearer <insert_the_access_token_here>"
  -H "Content-Length:0"
```

Example Response

```json
{
  "vehicleDataReadoutRefresh": {
    "id": "exvs3033320330937090",
    "asyncStatus": "Complete",
    "messageTimestamp": "2019-11-05T11:49:09+0000",
    "vehicleId": "WDD111111CAR01000",
    "message": [
      {
        "message": "Request to refresh the vehicle data readout was successful"
      }
    ]
  }
}
```

## General Error Handling

The API has a consistent error handling concept with the following error codes.

| **HTTP Error Code** | **Description**                                                                                                       |
| ------------------- | --------------------------------------------------------------------------------------------------------------------- |
| 400                 | The request could not complete due to malformed syntax.                                                               |
| 401                 | Invalid or missing authorization in header.                                                                           |
| 402                 | Payment required.                                                                                                     |
| 403                 | Forbidden.                                                                                                            |
| 404                 | The requested resource was not found.                                                                                 |
| 406                 | The resource identified by the request is not capable of generating a response that matches the given Accept headers. |
| 429                 | Quota limit is exceeded.                                                                                              |
| 500                 | An error occurred on the server side, e.g. a required service did not provide a valid response.                       |
| 501                 | The server does not support the functionality required.                                                               |
| 503                 | The server is unable to service the request due to a temporary unavailability condition.                              |
| 505                 | The server does not, or refuses to support the protocol associated with the request.                                  |
