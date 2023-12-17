package models

type TrackingStatus string

const (
    Created           TrackingStatus = "Created"
    PackageReceived   TrackingStatus = "Package received"
    InTransit         TrackingStatus = "In transit"
    OutForDelivery    TrackingStatus = "Out for delivery"
    DeliveryAttempted TrackingStatus = "Delivery attempted"
    Delivered         TrackingStatus = "Delivered"
    ReturnedToSender  TrackingStatus = "Returned to sender"
    Exception         TrackingStatus = "Exception"
)