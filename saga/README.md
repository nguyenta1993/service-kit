# s14e-saga
Saga Execution Coordinator 

### Pub/Sub:
- We have a publisher and a subscriber that can publish/subscribe to a channel. 
    - Subsciber will recieved msg.Message in method ReceiveMessage in registerd channel by implemented interface msg.MessageReceiver
    - Publisher have many ways to publish a command (PublishCommand, PublishReply, PublishEntityEvent, PublishEvent, ...etc) that can support various scenario.
     - Message must implemented coorresponding interface e.g: message of used by PublishCommand must be a Command (need to implemented msg.Command)


### Create a table to store the saga instance first: 
 `CREATE TABLE saga_instances (
    saga_name      text        NOT NULL,
    saga_id        text        NOT NULL,
    saga_data_name text        NOT NULL,
    saga_data      bytea       NOT NULL,
    current_step   int         NOT NULL,
    end_state      boolean     NOT NULL,
    compensating   boolean     NOT NULL,
    modified_at    timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (saga_name, saga_id)
)`

### To register a saga: 
 - Create a concrete SagaData (E.g UserSagaData) implemented core.SagaData interface and Register that SagaData to core `core.RegisterSagaData(UserSagaData{})`
 - Create a Saga Orchestrator (E.g UserSagaOrchestrator) in here we will define saga step.
    - Remote Step: Step that will publish a command to another service and will handle the result via Kafka
        - Other service need to return msg.Reply with either success or failure
    - Local Step: Step that will be executed in local
        - Action: Action will be called when saga start.
        - Compensation: Will called to compensate if the action fail to rollback.
- To start a saga, call the `Start` Method of the saga and the orchestrator will start.
- 