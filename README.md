# Heisprosjekt

## Datastrukturer
- State_matrix er et map, der key er siste tre tall i ip-adressen og verdiene til hver key er av typen Elevator_states
- Internal_order er et map som inneholder de interne knappetrykkene
- External_order er en matrise med opp- og nedknapp til hver etasje
- Elevator_states er en struct med ulike tilstander og verdier heisen kan ha
- Order_type er en struct med etasje og knappetrykk

## Moduler
- main
- driver: her defineres konstanter, og alt som har med hardware håndteres
- types: ulike struct defineres
- elevator: her er initialiseringsfunksjonen og funksjonen som kjører heisen. Alt som har med det fysiske på heisen settes her.
- orders: Her oppdateres verdier i State_matrix
- Network: dette er Anders Rønning Petersen sin network-modul, med noen modifikasjoner. I network.go printes Peer-oppdateringer, som er en oversikt over andre heiser som sender peer-oppdatering på samme kanal. Netverksprotokoll: UDP

## Antagelser:
- Antar at alle går inn i heisen når heisen stopper i en etasje. 
