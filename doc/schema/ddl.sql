CREATE TABLE public.app_app (
                                app_id bigserial NOT NULL,
                                app_name varchar(400) NOT NULL,
                                app_status smallint NOT NULL,
                                time_created bigint NOT NULL,
                                time_updated bigint NOT NULL,
                                CONSTRAINT app_app_pk PRIMARY KEY (app_id)
);
CREATE UNIQUE INDEX app_app_app_name_idx ON public.app_app (app_name);

CREATE TABLE public.app_token (
                                  token_id bigserial NOT NULL,
                                  app_id bigint NOT NULL,
                                  "token" varchar(400) NOT NULL,
                                  security_value varchar(400) NULL,
                                  token_type smallint NOT NULL,
                                  token_status smallint NOT NULL,
                                  expirydate_millis bigint NULL,
                                  time_created bigint NOT NULL,
                                  time_updated bigint NOT NULL,
                                  CONSTRAINT app_token_pk PRIMARY KEY (token_id)
);
CREATE UNIQUE INDEX app_token_token_idx ON public.app_token ("token");

INSERT INTO app_app(APP_NAME, APP_STATUS, TIME_CREATED, TIME_UPDATED) VALUES ('user', 0, 1, 1);

INSERT INTO app_token(app_id, token, token_type, token_status, expirydate_millis, time_created, time_updated)
    VALUES (1, 'userservtoken', 1, 0, 253373578095000, 1, 1);
INSERT INTO app_token(app_id, token, security_value, token_type, token_status, expirydate_millis, time_created, time_updated)
    VALUES (1, 'userservtokensso', 'userservtokenssovalue', 1, 0, 253373578095000, 1, 1);

CREATE TABLE public.user_user (
                                  user_id bigserial NOT NULL,
                                  user_type smallint NOT NULL,
                                  user_status smallint NOT NULL,
                                  user_tag1 bigint NOT NULL DEFAULT 0,
                                  user_tag2 bigint NOT NULL DEFAULT 0,
                                  user_name varchar(400) NULL,
                                  display_name varchar(400) NULL,
                                  time_created bigint NOT NULL,
                                  time_updated bigint NOT NULL,
                                  CONSTRAINT user_user_pk PRIMARY KEY (user_id)
);

CREATE TABLE public.user_credential (
                                        credential_id bigserial NOT NULL,
                                        user_id bigint NOT NULL,
                                        credential_type smallint NOT NULL,
                                        credential_key varchar(400) NOT NULL,
                                        credential_value varchar(400) NULL,
                                        credential_status smallint NOT NULL,
                                        time_created bigint NOT NULL,
                                        time_updated bigint NOT NULL,
                                        CONSTRAINT user_credential_pk PRIMARY KEY (credential_id)
);
CREATE UNIQUE INDEX user_credential_credential_key_idx ON public.user_credential (credential_key,credential_type);

CREATE TABLE public.sso_app (
                                sso_app_id bigserial NOT NULL,
                                app_id bigint NOT NULL,
                                sso_app_status smallint NOT NULL,
                                redirect_uri varchar(400) NOT NULL,
                                callback_uri varchar(400) NOT NULL,
                                time_created bigint NOT NULL,
                                time_updated bigint NOT NULL,
                                CONSTRAINT sso_app_pk PRIMARY KEY (sso_app_id)
);
CREATE INDEX sso_app_app_id_idx ON public.sso_app (app_id);

INSERT INTO sso_app(app_id, sso_app_status, redirect_uri, callback_uri, time_created, time_updated)
VALUES (1, 0, '*', 'http://example.com', 1, 1);


CREATE TABLE public.user_access_token (
                                          user_access_token_id bigserial NOT NULL,
                                          app_token_id bigint NOT NULL,
                                          user_id bigint NOT NULL,
                                          access_token varchar(4000) NOT NULL,
                                          expirydate_sec bigint NOT NULL,
                                          status smallint NOT NULL,
                                          time_created bigint NOT NULL,
                                          time_updated bigint NOT NULL,
                                          CONSTRAINT user_access_token_pk PRIMARY KEY (user_access_token_id)
);
CREATE INDEX user_access_token_app_token_id_idx ON public.user_access_token (app_token_id);
CREATE INDEX user_access_token_user_id_idx ON public.user_access_token (user_id);
CREATE UNIQUE INDEX user_access_token_access_token_idx ON public.user_access_token (access_token);
