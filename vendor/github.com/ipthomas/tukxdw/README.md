# tukxdw
Implementation of IHE XDW Actors
  XDW Content Creator
  XDW Content Consumer
  XDW Content Updator

Requires external Services:-
  IHE DSUB Broker WSE
  IHE PIXm or PIXv3 or PDQv3 server WSE
  MySql DB - connection made using DSN or if url provided by url.

The IHE XDW profile utilises the WS OASIS Human Task standard for the definition of Cross Document Workflow owners, tasks, inputs, outputs, completion behaviours, etc.

XDW requires a workflow definition for each workflow. Workflow definitions are defined using the IHE XDW profile / WS Human Tasks standards. When you use tukxdw to register a XDW definition with the TUK Event Service, the tukxdw transaction creates DSUB Broker Subscriptions for each XDW input and output task that has a type of '$XDSDocumentEntryTypeCode' in the XDW definition. The resulting broker subscription reference, NHS ID, XDW pathway, topic and expression for each subscription is persisted in the tuk event 'subscriptions' DB table. This enables received notifications from a DSUB Broker to be matched to a specific pathway and a specific task in that pathway. 

An event for each DSUB notification recieved is created and persisted with the event service. For an example implementation of a DSUB Broker Event Consumer that receives DSUB Broker Notify messages, parses IHE DSUB Notify message and persists the meta data to the TUK Event Service database table 'events', refer to github.com/ipthomas/tukdsub for local deployment and github/ipthomas/tukdsub_lambda for AWS deployment.

To query/update the TUK Event Service Database directly import github.com/ipthomas/tukdbint into your Go project
