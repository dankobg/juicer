--
-- PostgreSQL database cluster dump
--

SET default_transaction_read_only = off;

SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;

--
-- Drop databases (except postgres and template1)
--

DROP DATABASE test;




--
-- Drop roles
--

DROP ROLE test;


--
-- Roles
--

CREATE ROLE test;
ALTER ROLE test WITH SUPERUSER INHERIT CREATEROLE CREATEDB LOGIN REPLICATION BYPASSRLS PASSWORD 'SCRAM-SHA-256$4096:beNpk4LJObBGQfhLpc6ftQ==$Uzwf49nX/4STFdEUSXNZb1XVSvBW5VYZksPLTqp+E1w=:7nvG1645fgGFbzxsdWYBjKLMz//zYfAV8g20x52xeGA=';

--
-- User Configurations
--








--
-- Databases
--

--
-- Database "template1" dump
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3
-- Dumped by pg_dump version 16.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

UPDATE pg_catalog.pg_database SET datistemplate = false WHERE datname = 'template1';
DROP DATABASE template1;
--
-- Name: template1; Type: DATABASE; Schema: -; Owner: test
--

CREATE DATABASE template1 WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


ALTER DATABASE template1 OWNER TO test;

\connect template1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: DATABASE template1; Type: COMMENT; Schema: -; Owner: test
--

COMMENT ON DATABASE template1 IS 'default template for new databases';


--
-- Name: template1; Type: DATABASE PROPERTIES; Schema: -; Owner: test
--

ALTER DATABASE template1 IS_TEMPLATE = true;


\connect template1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: DATABASE template1; Type: ACL; Schema: -; Owner: test
--

REVOKE CONNECT,TEMPORARY ON DATABASE template1 FROM PUBLIC;
GRANT CONNECT ON DATABASE template1 TO PUBLIC;


--
-- PostgreSQL database dump complete
--

--
-- Database "postgres" dump
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3
-- Dumped by pg_dump version 16.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

DROP DATABASE postgres;
--
-- Name: postgres; Type: DATABASE; Schema: -; Owner: test
--

CREATE DATABASE postgres WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


ALTER DATABASE postgres OWNER TO test;

\connect postgres

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: DATABASE postgres; Type: COMMENT; Schema: -; Owner: test
--

COMMENT ON DATABASE postgres IS 'default administrative connection database';


--
-- PostgreSQL database dump complete
--

--
-- Database "test" dump
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3
-- Dumped by pg_dump version 16.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: test; Type: DATABASE; Schema: -; Owner: test
--

CREATE DATABASE test WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


ALTER DATABASE test OWNER TO test;

\connect test

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: btree_gin; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS btree_gin WITH SCHEMA public;


--
-- Name: EXTENSION btree_gin; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION btree_gin IS 'support for indexing common datatypes in GIN';


--
-- Name: pg_trgm; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;


--
-- Name: EXTENSION pg_trgm; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pg_trgm IS 'text similarity measurement and index searching based on trigrams';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: continuity_containers; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.continuity_containers (
    id uuid NOT NULL,
    identity_id uuid,
    name character varying(255) NOT NULL,
    payload jsonb,
    expires_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    nid uuid
);


ALTER TABLE public.continuity_containers OWNER TO test;

--
-- Name: courier_message_dispatches; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.courier_message_dispatches (
    id uuid NOT NULL,
    message_id uuid NOT NULL,
    status character varying(7) NOT NULL,
    error json,
    nid uuid NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.courier_message_dispatches OWNER TO test;

--
-- Name: courier_messages; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.courier_messages (
    id uuid NOT NULL,
    type integer NOT NULL,
    status integer NOT NULL,
    body text NOT NULL,
    subject character varying(255) NOT NULL,
    recipient character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    template_type character varying(255) DEFAULT ''::character varying NOT NULL,
    template_data bytea,
    nid uuid,
    send_count integer DEFAULT 0 NOT NULL,
    channel character varying(32)
);


ALTER TABLE public.courier_messages OWNER TO test;

--
-- Name: identities; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identities (
    id uuid NOT NULL,
    schema_id character varying(2048) NOT NULL,
    traits jsonb NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    nid uuid,
    state character varying(255) DEFAULT 'active'::character varying NOT NULL,
    state_changed_at timestamp without time zone,
    metadata_public jsonb,
    metadata_admin jsonb,
    available_aal character varying(4),
    organization_id uuid
);


ALTER TABLE public.identities OWNER TO test;

--
-- Name: identity_credential_identifiers; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_credential_identifiers (
    id uuid NOT NULL,
    identifier character varying(255) NOT NULL,
    identity_credential_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    nid uuid,
    identity_credential_type_id uuid NOT NULL
);


ALTER TABLE public.identity_credential_identifiers OWNER TO test;

--
-- Name: identity_credential_types; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_credential_types (
    id uuid NOT NULL,
    name character varying(32) NOT NULL
);


ALTER TABLE public.identity_credential_types OWNER TO test;

--
-- Name: identity_credentials; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_credentials (
    id uuid NOT NULL,
    config jsonb NOT NULL,
    identity_credential_type_id uuid NOT NULL,
    identity_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    nid uuid,
    version integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.identity_credentials OWNER TO test;

--
-- Name: identity_login_codes; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_login_codes (
    id uuid NOT NULL,
    code character varying(64) NOT NULL,
    address character varying(255) NOT NULL,
    address_type character(36) NOT NULL,
    used_at timestamp without time zone,
    expires_at timestamp without time zone DEFAULT '2000-01-01 00:00:00'::timestamp without time zone NOT NULL,
    issued_at timestamp without time zone DEFAULT '2000-01-01 00:00:00'::timestamp without time zone NOT NULL,
    selfservice_login_flow_id uuid NOT NULL,
    identity_id uuid NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    nid uuid NOT NULL
);


ALTER TABLE public.identity_login_codes OWNER TO test;

--
-- Name: identity_recovery_addresses; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_recovery_addresses (
    id uuid NOT NULL,
    via character varying(16) NOT NULL,
    value character varying(400) NOT NULL,
    identity_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    nid uuid
);


ALTER TABLE public.identity_recovery_addresses OWNER TO test;

--
-- Name: identity_recovery_codes; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_recovery_codes (
    id uuid NOT NULL,
    code character varying(64) NOT NULL,
    used_at timestamp without time zone,
    identity_recovery_address_id uuid,
    code_type integer NOT NULL,
    expires_at timestamp without time zone DEFAULT '2000-01-01 00:00:00'::timestamp without time zone NOT NULL,
    issued_at timestamp without time zone DEFAULT '2000-01-01 00:00:00'::timestamp without time zone NOT NULL,
    selfservice_recovery_flow_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    nid uuid NOT NULL,
    identity_id uuid NOT NULL
);


ALTER TABLE public.identity_recovery_codes OWNER TO test;

--
-- Name: identity_recovery_tokens; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_recovery_tokens (
    id uuid NOT NULL,
    token character varying(64) NOT NULL,
    used boolean DEFAULT false NOT NULL,
    used_at timestamp without time zone,
    identity_recovery_address_id uuid,
    selfservice_recovery_flow_id uuid,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    expires_at timestamp without time zone DEFAULT '2000-01-01 00:00:00'::timestamp without time zone NOT NULL,
    issued_at timestamp without time zone DEFAULT '2000-01-01 00:00:00'::timestamp without time zone NOT NULL,
    nid uuid,
    identity_id uuid NOT NULL,
    token_type integer DEFAULT 0 NOT NULL,
    CONSTRAINT identity_recovery_tokens_token_type_ck CHECK (((token_type = 1) OR (token_type = 2)))
);


ALTER TABLE public.identity_recovery_tokens OWNER TO test;

--
-- Name: identity_registration_codes; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_registration_codes (
    id uuid NOT NULL,
    code character varying(64) NOT NULL,
    address character varying(255) NOT NULL,
    address_type character(36) NOT NULL,
    used_at timestamp without time zone,
    expires_at timestamp without time zone DEFAULT '2000-01-01 00:00:00'::timestamp without time zone NOT NULL,
    issued_at timestamp without time zone DEFAULT '2000-01-01 00:00:00'::timestamp without time zone NOT NULL,
    selfservice_registration_flow_id uuid NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    nid uuid NOT NULL
);


ALTER TABLE public.identity_registration_codes OWNER TO test;

--
-- Name: identity_verifiable_addresses; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_verifiable_addresses (
    id uuid NOT NULL,
    status character varying(16) NOT NULL,
    via character varying(16) NOT NULL,
    verified boolean NOT NULL,
    value character varying(400) NOT NULL,
    verified_at timestamp without time zone,
    identity_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    nid uuid
);


ALTER TABLE public.identity_verifiable_addresses OWNER TO test;

--
-- Name: identity_verification_codes; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_verification_codes (
    id uuid NOT NULL,
    code_hmac character varying(64) NOT NULL,
    used_at timestamp without time zone,
    identity_verifiable_address_id uuid,
    expires_at timestamp without time zone DEFAULT '2000-01-01 00:00:00'::timestamp without time zone NOT NULL,
    issued_at timestamp without time zone DEFAULT '2000-01-01 00:00:00'::timestamp without time zone NOT NULL,
    selfservice_verification_flow_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    nid uuid NOT NULL
);


ALTER TABLE public.identity_verification_codes OWNER TO test;

--
-- Name: identity_verification_tokens; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.identity_verification_tokens (
    id uuid NOT NULL,
    token character varying(64) NOT NULL,
    used boolean DEFAULT false NOT NULL,
    used_at timestamp without time zone,
    expires_at timestamp without time zone NOT NULL,
    issued_at timestamp without time zone NOT NULL,
    identity_verifiable_address_id uuid NOT NULL,
    selfservice_verification_flow_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    nid uuid
);


ALTER TABLE public.identity_verification_tokens OWNER TO test;

--
-- Name: keto_relation_tuples; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.keto_relation_tuples (
    shard_id uuid NOT NULL,
    nid uuid NOT NULL,
    namespace character varying(200) NOT NULL,
    object uuid NOT NULL,
    relation character varying(64) NOT NULL,
    subject_id uuid,
    subject_set_namespace character varying(200),
    subject_set_object uuid,
    subject_set_relation character varying(64),
    commit_time timestamp without time zone NOT NULL,
    CONSTRAINT chk_keto_rt_uuid_subject_type CHECK ((((subject_id IS NULL) AND (subject_set_namespace IS NOT NULL) AND (subject_set_object IS NOT NULL) AND (subject_set_relation IS NOT NULL)) OR ((subject_id IS NOT NULL) AND (subject_set_namespace IS NULL) AND (subject_set_object IS NULL) AND (subject_set_relation IS NULL))))
);


ALTER TABLE public.keto_relation_tuples OWNER TO test;

--
-- Name: keto_uuid_mappings; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.keto_uuid_mappings (
    id uuid NOT NULL,
    string_representation text NOT NULL
);


ALTER TABLE public.keto_uuid_mappings OWNER TO test;

--
-- Name: networks; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.networks (
    id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.networks OWNER TO test;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.schema_migration (
    version character varying(48) NOT NULL,
    version_self integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO test;

--
-- Name: selfservice_errors; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.selfservice_errors (
    id uuid NOT NULL,
    errors jsonb NOT NULL,
    seen_at timestamp without time zone,
    was_seen boolean NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    csrf_token character varying(255) DEFAULT ''::character varying NOT NULL,
    nid uuid
);


ALTER TABLE public.selfservice_errors OWNER TO test;

--
-- Name: selfservice_login_flows; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.selfservice_login_flows (
    id uuid NOT NULL,
    request_url text NOT NULL,
    issued_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    active_method character varying(32) NOT NULL,
    csrf_token character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    forced boolean DEFAULT false NOT NULL,
    type character varying(16) DEFAULT 'browser'::character varying NOT NULL,
    ui jsonb,
    nid uuid,
    requested_aal character varying(4) DEFAULT 'aal1'::character varying NOT NULL,
    internal_context jsonb NOT NULL,
    oauth2_login_challenge uuid,
    oauth2_login_challenge_data text,
    state character varying(255),
    submit_count integer DEFAULT 0 NOT NULL,
    organization_id uuid
);


ALTER TABLE public.selfservice_login_flows OWNER TO test;

--
-- Name: selfservice_recovery_flows; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.selfservice_recovery_flows (
    id uuid NOT NULL,
    request_url text NOT NULL,
    issued_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    active_method character varying(32),
    csrf_token character varying(255) NOT NULL,
    state character varying(32) NOT NULL,
    recovered_identity_id uuid,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    type character varying(16) DEFAULT 'browser'::character varying NOT NULL,
    ui jsonb,
    nid uuid,
    submit_count integer DEFAULT 0 NOT NULL,
    skip_csrf_check boolean DEFAULT false NOT NULL
);


ALTER TABLE public.selfservice_recovery_flows OWNER TO test;

--
-- Name: selfservice_registration_flows; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.selfservice_registration_flows (
    id uuid NOT NULL,
    request_url text NOT NULL,
    issued_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    active_method character varying(32) NOT NULL,
    csrf_token character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    type character varying(16) DEFAULT 'browser'::character varying NOT NULL,
    ui jsonb,
    nid uuid,
    internal_context jsonb NOT NULL,
    oauth2_login_challenge uuid,
    oauth2_login_challenge_data text,
    state character varying(255),
    submit_count integer DEFAULT 0 NOT NULL,
    organization_id uuid
);


ALTER TABLE public.selfservice_registration_flows OWNER TO test;

--
-- Name: selfservice_settings_flows; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.selfservice_settings_flows (
    id uuid NOT NULL,
    request_url text NOT NULL,
    issued_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    identity_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    active_method character varying(32),
    state character varying(255) DEFAULT 'show_form'::character varying NOT NULL,
    type character varying(16) DEFAULT 'browser'::character varying NOT NULL,
    ui jsonb,
    nid uuid,
    internal_context jsonb NOT NULL
);


ALTER TABLE public.selfservice_settings_flows OWNER TO test;

--
-- Name: selfservice_verification_flows; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.selfservice_verification_flows (
    id uuid NOT NULL,
    request_url text NOT NULL,
    issued_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    csrf_token character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    type character varying(16) DEFAULT 'browser'::character varying NOT NULL,
    state character varying(255) DEFAULT 'show_form'::character varying NOT NULL,
    active_method character varying(32),
    ui jsonb,
    nid uuid,
    submit_count integer DEFAULT 0 NOT NULL,
    oauth2_login_challenge text,
    session_id uuid,
    identity_id uuid,
    authentication_methods json
);


ALTER TABLE public.selfservice_verification_flows OWNER TO test;

--
-- Name: session_devices; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.session_devices (
    id uuid NOT NULL,
    ip_address character varying(50) DEFAULT ''::character varying,
    user_agent character varying(512) DEFAULT ''::character varying,
    location character varying(512) DEFAULT ''::character varying,
    nid uuid NOT NULL,
    session_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.session_devices OWNER TO test;

--
-- Name: session_token_exchanges; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.session_token_exchanges (
    id uuid NOT NULL,
    nid uuid NOT NULL,
    flow_id uuid NOT NULL,
    session_id uuid,
    init_code character varying(64) NOT NULL,
    return_to_code character varying(64) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.session_token_exchanges OWNER TO test;

--
-- Name: sessions; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.sessions (
    id uuid NOT NULL,
    issued_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    authenticated_at timestamp without time zone NOT NULL,
    identity_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    token character varying(39),
    active boolean DEFAULT false,
    nid uuid,
    logout_token character varying(39),
    aal character varying(4) DEFAULT 'aal1'::character varying NOT NULL,
    authentication_methods jsonb NOT NULL
);


ALTER TABLE public.sessions OWNER TO test;

--
-- Data for Name: continuity_containers; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.continuity_containers (id, identity_id, name, payload, expires_at, created_at, updated_at, nid) FROM stdin;
\.


--
-- Data for Name: courier_message_dispatches; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.courier_message_dispatches (id, message_id, status, error, nid, created_at, updated_at) FROM stdin;
1c32dc9a-8d98-4153-8824-f52742f3a353	90b4378f-947e-4920-abbd-f799586a57f9	success	null	e796004f-1c18-4960-9388-dcd5835741db	2024-05-18 11:25:46.965435	2024-05-18 11:25:46.965435
\.


--
-- Data for Name: courier_messages; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.courier_messages (id, type, status, body, subject, recipient, created_at, updated_at, template_type, template_data, nid, send_count, channel) FROM stdin;
90b4378f-947e-4920-abbd-f799586a57f9	1	2	Hi,\n\nplease verify your account by entering the following code:\n\n395337\n\nor clicking the following link:\n\nhttps://juicer-dev.xyz/kratos/self-service/verification?code=395337&flow=f38bd08b-a3ee-475b-8a24-eb89350a3b3f\n	Please verify your email address	test@test.com	2024-05-18 11:25:45.998977	2024-05-18 11:25:45.998977	verification_code_valid	\\x7b22746f223a227465737440746573742e636f6d222c22766572696669636174696f6e5f75726c223a2268747470733a2f2f6a75696365722d6465762e78797a2f6b7261746f732f73656c662d736572766963652f766572696669636174696f6e3f636f64653d3339353333375c7530303236666c6f773d66333862643038622d613365652d343735622d386132342d656238393335306133623366222c22766572696669636174696f6e5f636f6465223a22333935333337222c226964656e74697479223a7b22637265617465645f6174223a22323032342d30352d31385431313a32353a34352e3937323339335a222c226964223a2262623431356432332d643430612d343035642d383237392d313066616232633664396130222c226d657461646174615f7075626c6963223a6e756c6c2c226f7267616e697a6174696f6e5f6964223a6e756c6c2c227265636f766572795f616464726573736573223a5b7b22637265617465645f6174223a22323032342d30352d31385431313a32353a34352e3937353439355a222c226964223a2261656561303338352d383630662d343235392d393030642d613739623231363864323233222c22757064617465645f6174223a22323032342d30352d31385431313a32353a34352e3937353439355a222c2276616c7565223a227465737440746573742e636f6d222c22766961223a22656d61696c227d5d2c22736368656d615f6964223a22637573746f6d6572222c22736368656d615f75726c223a2268747470733a2f2f6a75696365722d6465762e78797a2f6b7261746f732f736368656d61732f5933567a644739745a5849222c227374617465223a22616374697665222c2273746174655f6368616e6765645f6174223a22323032342d30352d31385431313a32353a34352e3937303733373237355a222c22747261697473223a7b226176617461725f75726c223a22222c22656d61696c223a227465737440746573742e636f6d222c2266697273745f6e616d65223a22546573744669727374222c226c6173745f6e616d65223a22546573744c617374227d2c22757064617465645f6174223a22323032342d30352d31385431313a32353a34352e3937323339335a222c2276657269666961626c655f616464726573736573223a5b7b22637265617465645f6174223a22323032342d30352d31385431313a32353a34352e3937333735335a222c226964223a2235343336363833312d343966362d343863332d383739632d613631343738303661383637222c22737461747573223a2270656e64696e67222c22757064617465645f6174223a22323032342d30352d31385431313a32353a34352e3937333735335a222c2276616c7565223a227465737440746573742e636f6d222c227665726966696564223a66616c73652c22766961223a22656d61696c227d5d7d2c22726571756573745f75726c223a2268747470733a2f2f6a75696365722d6465762e78797a2f73656c662d736572766963652f726567697374726174696f6e2f62726f77736572227d	e796004f-1c18-4960-9388-dcd5835741db	1	email
\.


--
-- Data for Name: identities; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identities (id, schema_id, traits, created_at, updated_at, nid, state, state_changed_at, metadata_public, metadata_admin, available_aal, organization_id) FROM stdin;
bb415d23-d40a-405d-8279-10fab2c6d9a0	customer	{"email": "test@test.com", "last_name": "TestLast", "avatar_url": "", "first_name": "TestFirst"}	2024-05-18 11:25:45.972393	2024-05-18 11:25:45.972393	e796004f-1c18-4960-9388-dcd5835741db	active	2024-05-18 11:25:45.970737	\N	\N	aal1	\N
\.


--
-- Data for Name: identity_credential_identifiers; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_credential_identifiers (id, identifier, identity_credential_id, created_at, updated_at, nid, identity_credential_type_id) FROM stdin;
e2f7aade-5ac5-4986-9623-687995d849d7	test@test.com	ac552fa8-1118-40c8-9989-7aff8a84295b	2024-05-18 11:25:45.980547	2024-05-18 11:25:45.980547	e796004f-1c18-4960-9388-dcd5835741db	78c1b41d-8341-4507-aa60-aff1d4369670
\.


--
-- Data for Name: identity_credential_types; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_credential_types (id, name) FROM stdin;
78c1b41d-8341-4507-aa60-aff1d4369670	password
6fa5e2e0-bfce-4631-b62b-cf2b0252b289	oidc
5e29b036-aa47-457f-9fe6-aa8b854a752b	totp
567a0730-7f48-4dd7-a13d-df87a51c245f	lookup_secret
6b213fa0-e6ad-46cb-8878-b088d2ce2e3c	webauthn
14f3b7e2-8725-4068-be39-8a796485fd97	code
\.


--
-- Data for Name: identity_credentials; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_credentials (id, config, identity_credential_type_id, identity_id, created_at, updated_at, nid, version) FROM stdin;
ac552fa8-1118-40c8-9989-7aff8a84295b	{"hashed_password": "$2a$12$JKKFuTFMnHZd30f99aopcOCcs7Ntd9k0H22XhUHmC4ZZd8elRZ5Wi"}	78c1b41d-8341-4507-aa60-aff1d4369670	bb415d23-d40a-405d-8279-10fab2c6d9a0	2024-05-18 11:25:45.979103	2024-05-18 11:25:45.979103	e796004f-1c18-4960-9388-dcd5835741db	0
\.


--
-- Data for Name: identity_login_codes; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_login_codes (id, code, address, address_type, used_at, expires_at, issued_at, selfservice_login_flow_id, identity_id, created_at, updated_at, nid) FROM stdin;
\.


--
-- Data for Name: identity_recovery_addresses; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_recovery_addresses (id, via, value, identity_id, created_at, updated_at, nid) FROM stdin;
aeea0385-860f-4259-900d-a79b2168d223	email	test@test.com	bb415d23-d40a-405d-8279-10fab2c6d9a0	2024-05-18 11:25:45.975495	2024-05-18 11:25:45.975495	e796004f-1c18-4960-9388-dcd5835741db
\.


--
-- Data for Name: identity_recovery_codes; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_recovery_codes (id, code, used_at, identity_recovery_address_id, code_type, expires_at, issued_at, selfservice_recovery_flow_id, created_at, updated_at, nid, identity_id) FROM stdin;
\.


--
-- Data for Name: identity_recovery_tokens; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_recovery_tokens (id, token, used, used_at, identity_recovery_address_id, selfservice_recovery_flow_id, created_at, updated_at, expires_at, issued_at, nid, identity_id, token_type) FROM stdin;
\.


--
-- Data for Name: identity_registration_codes; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_registration_codes (id, code, address, address_type, used_at, expires_at, issued_at, selfservice_registration_flow_id, created_at, updated_at, nid) FROM stdin;
\.


--
-- Data for Name: identity_verifiable_addresses; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_verifiable_addresses (id, status, via, verified, value, verified_at, identity_id, created_at, updated_at, nid) FROM stdin;
54366831-49f6-48c3-879c-a6147806a867	completed	email	t	test@test.com	2024-05-18 11:25:58.638391	bb415d23-d40a-405d-8279-10fab2c6d9a0	2024-05-18 11:25:45.973753	2024-05-18 11:25:45.973753	e796004f-1c18-4960-9388-dcd5835741db
\.


--
-- Data for Name: identity_verification_codes; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_verification_codes (id, code_hmac, used_at, identity_verifiable_address_id, expires_at, issued_at, selfservice_verification_flow_id, created_at, updated_at, nid) FROM stdin;
d5b640b6-8c94-4d78-b61e-6fda67731394	172f22b6837ad69506716377c913e3e7cab662d849053f43f86319d5d01c066c	2024-05-18 11:25:58.627061	54366831-49f6-48c3-879c-a6147806a867	2024-05-18 11:40:45.994271	2024-05-18 11:25:45.994271	f38bd08b-a3ee-475b-8a24-eb89350a3b3f	2024-05-18 11:25:45.994332	2024-05-18 11:25:45.994332	e796004f-1c18-4960-9388-dcd5835741db
\.


--
-- Data for Name: identity_verification_tokens; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.identity_verification_tokens (id, token, used, used_at, expires_at, issued_at, identity_verifiable_address_id, selfservice_verification_flow_id, created_at, updated_at, nid) FROM stdin;
\.


--
-- Data for Name: keto_relation_tuples; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.keto_relation_tuples (shard_id, nid, namespace, object, relation, subject_id, subject_set_namespace, subject_set_object, subject_set_relation, commit_time) FROM stdin;
\.


--
-- Data for Name: keto_uuid_mappings; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.keto_uuid_mappings (id, string_representation) FROM stdin;
\.


--
-- Data for Name: networks; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.networks (id, created_at, updated_at) FROM stdin;
e796004f-1c18-4960-9388-dcd5835741db	2024-05-18 11:23:08.822098	2024-05-18 11:23:08.822098
\.


--
-- Data for Name: schema_migration; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.schema_migration (version, version_self) FROM stdin;
20150100000001000000	0
20201110175414000000	0
20201110175414000001	0
20210623162417000000	0
20210623162417000001	0
20210623162417000002	0
20210623162417000003	0
20210914134624000000	0
20220217152313000000	0
20220512151000000000	0
20220513200300000000	0
20220513200400000000	0
20220513200400000001	0
20220513200500000000	0
20220513200600000000	0
20220513200600000001	0
20230228091200000000	0
20191100000001000000	0
20191100000001000001	0
20191100000001000002	0
20191100000001000003	0
20191100000001000004	0
20191100000001000005	0
20191100000002000000	0
20191100000002000001	0
20191100000002000002	0
20191100000002000003	0
20191100000002000004	0
20191100000003000000	0
20191100000004000000	0
20191100000006000000	0
20191100000007000000	0
20191100000008000000	0
20191100000008000001	0
20191100000008000002	0
20191100000008000003	0
20191100000008000004	0
20191100000008000005	0
20191100000010000000	0
20191100000010000001	0
20191100000011000000	0
20191100000012000000	0
20200317160354000000	0
20200317160354000001	0
20200317160354000002	0
20200317160354000003	0
20200317160354000004	0
20200401183443000000	0
20200402142539000000	0
20200402142539000001	0
20200402142539000002	0
20200519101057000000	0
20200519101057000001	0
20200519101057000002	0
20200519101057000003	0
20200519101057000004	0
20200519101057000005	0
20200519101057000006	0
20200519101057000007	0
20200601101000000000	0
20200605111551000000	0
20200605111551000001	0
20200605111551000002	0
20200607165100000000	0
20200607165100000001	0
20200705105359000000	0
20200810141652000000	0
20200810141652000001	0
20200810141652000002	0
20200810141652000003	0
20200810141652000004	0
20200810161022000000	0
20200810161022000001	0
20200810161022000002	0
20200810161022000003	0
20200810161022000004	0
20200810161022000005	0
20200810161022000006	0
20200810161022000007	0
20200810161022000008	0
20200810162450000000	0
20200810162450000001	0
20200810162450000002	0
20200810162450000003	0
20200812124254000000	0
20200812124254000001	0
20200812124254000002	0
20200812124254000003	0
20200812124254000004	0
20200812160551000000	0
20200830121710000000	0
20200830130642000000	0
20200830130642000001	0
20200830130642000002	0
20200830130642000003	0
20200830130642000004	0
20200830130642000005	0
20200830130642000006	0
20200830130642000007	0
20200830130643000000	0
20200830130644000000	0
20200830130644000001	0
20200830130645000000	0
20200830130646000000	0
20200830130646000001	0
20200830130646000002	0
20200830154602000000	0
20200830154602000001	0
20200830154602000002	0
20200830154602000003	0
20200830154602000004	0
20200830172221000000	0
20200830172221000001	0
20200830172221000002	0
20200830172221000003	0
20200831110752000000	0
20200831110752000001	0
20200831110752000002	0
20200831110752000003	0
20200831110752000004	0
20200831110752000005	0
20200831110752000006	0
20200831110752000007	0
20201201161451000000	0
20201201161451000001	0
20210307130558000000	0
20210307130559000000	0
20210307130559000001	0
20210311102338000000	0
20210311102338000001	0
20210311102338000002	0
20210311102338000003	0
20210311102338000004	0
20210311102338000005	0
20210311102338000006	0
20210311102338000007	0
20210311102338000008	0
20210311102338000009	0
20210311102338000010	0
20210311102338000011	0
20210311102338000012	0
20210311102338000013	0
20210311102338000014	0
20210311102338000015	0
20210311102338000016	0
20210311102338000017	0
20210311102338000018	0
20210311102338000019	0
20210311102338000020	0
20210311102338000021	0
20210311102338000022	0
20210311102338000023	0
20210311102338000024	0
20210410175418000000	0
20210410175418000001	0
20210410175418000002	0
20210410175418000003	0
20210410175418000004	0
20210410175418000005	0
20210410175418000006	0
20210410175418000007	0
20210410175418000008	0
20210410175418000009	0
20210410175418000010	0
20210410175418000011	0
20210410175418000012	0
20210410175418000013	0
20210410175418000014	0
20210410175418000015	0
20210410175418000016	0
20210410175418000017	0
20210410175418000018	0
20210410175418000019	0
20210410175418000020	0
20210410175418000021	0
20210410175418000022	0
20210410175418000023	0
20210410175418000024	0
20210410175418000025	0
20210410175418000026	0
20210410175418000027	0
20210410175418000028	0
20210410175418000029	0
20210410175418000030	0
20210410175418000031	0
20210410175418000032	0
20210410175418000033	0
20210410175418000034	0
20210410175418000035	0
20210410175418000036	0
20210410175418000037	0
20210410175418000038	0
20210410175418000039	0
20210410175418000040	0
20210410175418000041	0
20210410175418000042	0
20210410175418000043	0
20210410175418000044	0
20210410175418000045	0
20210410175418000046	0
20210410175418000047	0
20210410175418000048	0
20210410175418000049	0
20210410175418000050	0
20210410175418000051	0
20210410175418000052	0
20210410175418000053	0
20210410175418000054	0
20210410175418000055	0
20210410175418000056	0
20210410175418000057	0
20210410175418000058	0
20210410175418000059	0
20210410175418000060	0
20210410175418000061	0
20210410175418000062	0
20210410175418000063	0
20210410175418000064	0
20210410175418000065	0
20210410175418000066	0
20210410175418000067	0
20210410175418000068	0
20210410175418000069	0
20210410175418000070	0
20210410175418000071	0
20210410175418000072	0
20210410175418000073	0
20210410175418000074	0
20210410175418000075	0
20210410175418000076	0
20210410175418000077	0
20210410175418000078	0
20210410175418000079	0
20210410175418000080	0
20210410175418000081	0
20210410175418000082	0
20210410175418000083	0
20210410175418000084	0
20210410175418000085	0
20210410175418000086	0
20210410175418000087	0
20210410175418000088	0
20210410175418000089	0
20210504121624000000	0
20210504121624000001	0
20210618103120000000	0
20210618103120000001	0
20210618103120000002	0
20210618103120000003	0
20210618103120000004	0
20210805112414000000	0
20210805112414000001	0
20210805112414000002	0
20210805122535000000	0
20210810153530000000	0
20210810153530000001	0
20210810153530000002	0
20210810153530000003	0
20210810153530000004	0
20210813150152000000	0
20210816113956000000	0
20210816142650000000	0
20210816142650000001	0
20210816142650000002	0
20210816142650000003	0
20210816142650000004	0
20210816142650000005	0
20210817181232000000	0
20210817181232000001	0
20210817181232000002	0
20210817181232000003	0
20210817181232000004	0
20210817181232000005	0
20210829131458000000	0
20210913095309000000	0
20210913095309000001	0
20210913095309000002	0
20210913095309000003	0
20210913095309000004	0
20220118104539000000	0
20220118104539000001	0
20220118104539000002	0
20220118104539000003	0
20220301102701000000	0
20220301102701000001	0
20220420102701000000	0
20220512102703000000	0
20220607000001000000	0
20220610155809000000	0
20220802103909000000	0
20220824165300000000	0
20220824165300000001	0
20220824165300000002	0
20220825134336000000	0
20220825134336000001	0
20220901123209000000	0
20220907132836000000	0
20221024182336000000	0
20221205092803000000	0
20221214101328000000	0
20221220124639000000	0
20230104193739000000	0
20230216142104000000	0
20230313141439000000	0
20230313141439000001	0
20230322144139000001	0
20230405000000000001	0
20230614000001000000	0
20230619000000000001	0
20230626000000000001	0
20230703143600000001	0
20230705000000000001	0
20230706000000000001	0
20230707133700000000	0
20230707133700000001	0
20230712173852000000	0
20230811000000000001	0
20230818000000000001	0
20230823000000000001	0
20230907085000000000	0
20230920171028000000	0
20231130094628000000	0
20240119094628000000	0
\.


--
-- Data for Name: selfservice_errors; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.selfservice_errors (id, errors, seen_at, was_seen, created_at, updated_at, csrf_token, nid) FROM stdin;
\.


--
-- Data for Name: selfservice_login_flows; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.selfservice_login_flows (id, request_url, issued_at, expires_at, active_method, csrf_token, created_at, updated_at, forced, type, ui, nid, requested_aal, internal_context, oauth2_login_challenge, oauth2_login_challenge_data, state, submit_count, organization_id) FROM stdin;
3e2b092a-1c7e-45ac-a216-2cfd1d97ee3e	https://juicer-dev.xyz/self-service/login/browser	2024-05-18 11:24:05.614924	2024-05-18 11:34:05.614924		C1/xuVDGaKrQsyuibhaUKLALIYbKE4giKjpRscJ2Y4aGE001MO5VZRFwM0IBxJYkxOwUJVvejatZDKQ5Ad6ESg==	2024-05-18 11:24:05.617304	2024-05-18 11:24:05.617304	f	browser	{"nodes": [{"meta": {"label": {"id": 1010002, "text": "Sign in with discord", "type": "info", "context": {"provider": "discord"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "discord", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with facebook", "type": "info", "context": {"provider": "facebook"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "facebook", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with github", "type": "info", "context": {"provider": "github"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "github", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with google", "type": "info", "context": {"provider": "google"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "google", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with slack", "type": "info", "context": {"provider": "slack"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "slack", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with spotify", "type": "info", "context": {"provider": "spotify"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "spotify", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with twitch", "type": "info", "context": {"provider": "twitch"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "twitch", "disabled": false, "node_type": "input"}}, {"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "C1/xuVDGaKrQsyuibhaUKLALIYbKE4giKjpRscJ2Y4aGE001MO5VZRFwM0IBxJYkxOwUJVvejatZDKQ5Ad6ESg==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "E-Mail", "type": "info", "context": {"title": "E-Mail"}}}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "identifier", "type": "text", "value": "", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070001, "text": "Password", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "password", "type": "password", "disabled": false, "required": true, "node_type": "input", "autocomplete": "current-password"}}, {"meta": {"label": {"id": 1010001, "text": "Sign in", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "password", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/login?flow=3e2b092a-1c7e-45ac-a216-2cfd1d97ee3e", "method": "POST"}	e796004f-1c18-4960-9388-dcd5835741db	aal1	{}	\N	\N	choose_method	0	\N
0c941e3f-eef1-4080-97ef-b3096c0740c7	https://juicer-dev.xyz/self-service/login/browser	2024-05-18 11:24:28.310728	2024-05-18 11:34:28.310728		EklXw2PaT2C52XKg91Q24I8jMx1NVtMnNIKJ9cthcj2fBetPA/Jyr3gaakCYhjTs+8QGvtyb1q5HtHx9CMmV8Q==	2024-05-18 11:24:28.312672	2024-05-18 11:24:28.312672	f	browser	{"nodes": [{"meta": {"label": {"id": 1010002, "text": "Sign in with discord", "type": "info", "context": {"provider": "discord"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "discord", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with facebook", "type": "info", "context": {"provider": "facebook"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "facebook", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with github", "type": "info", "context": {"provider": "github"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "github", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with google", "type": "info", "context": {"provider": "google"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "google", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with slack", "type": "info", "context": {"provider": "slack"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "slack", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with spotify", "type": "info", "context": {"provider": "spotify"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "spotify", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with twitch", "type": "info", "context": {"provider": "twitch"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "twitch", "disabled": false, "node_type": "input"}}, {"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "EklXw2PaT2C52XKg91Q24I8jMx1NVtMnNIKJ9cthcj2fBetPA/Jyr3gaakCYhjTs+8QGvtyb1q5HtHx9CMmV8Q==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "E-Mail", "type": "info", "context": {"title": "E-Mail"}}}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "identifier", "type": "text", "value": "", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070001, "text": "Password", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "password", "type": "password", "disabled": false, "required": true, "node_type": "input", "autocomplete": "current-password"}}, {"meta": {"label": {"id": 1010001, "text": "Sign in", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "password", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/login?flow=0c941e3f-eef1-4080-97ef-b3096c0740c7", "method": "POST"}	e796004f-1c18-4960-9388-dcd5835741db	aal1	{}	\N	\N	choose_method	0	\N
ef1012ab-a33a-4d6b-908d-4f0f7ef4d194	https://juicer-dev.xyz/self-service/login/browser	2024-05-18 11:24:32.644001	2024-05-18 11:34:32.644001		SQXvztFFanlLsu35L2gt+IuG/z0FsUjEx39axKvkuVvESVNCsW1Xtopx9RlAui/0/2HKnpR8TU20Sa9MaExelw==	2024-05-18 11:24:32.649094	2024-05-18 11:24:32.649094	f	browser	{"nodes": [{"meta": {"label": {"id": 1010002, "text": "Sign in with discord", "type": "info", "context": {"provider": "discord"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "discord", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with facebook", "type": "info", "context": {"provider": "facebook"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "facebook", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with github", "type": "info", "context": {"provider": "github"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "github", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with google", "type": "info", "context": {"provider": "google"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "google", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with slack", "type": "info", "context": {"provider": "slack"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "slack", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with spotify", "type": "info", "context": {"provider": "spotify"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "spotify", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with twitch", "type": "info", "context": {"provider": "twitch"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "twitch", "disabled": false, "node_type": "input"}}, {"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "SQXvztFFanlLsu35L2gt+IuG/z0FsUjEx39axKvkuVvESVNCsW1Xtopx9RlAui/0/2HKnpR8TU20Sa9MaExelw==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "E-Mail", "type": "info", "context": {"title": "E-Mail"}}}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "identifier", "type": "text", "value": "", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070001, "text": "Password", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "password", "type": "password", "disabled": false, "required": true, "node_type": "input", "autocomplete": "current-password"}}, {"meta": {"label": {"id": 1010001, "text": "Sign in", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "password", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/login?flow=ef1012ab-a33a-4d6b-908d-4f0f7ef4d194", "method": "POST"}	e796004f-1c18-4960-9388-dcd5835741db	aal1	{}	\N	\N	choose_method	0	\N
753c9dc8-a04a-4655-8f9e-66cff204dc50	https://juicer-dev.xyz/self-service/login/browser	2024-05-18 11:25:04.465901	2024-05-18 11:35:04.465901		7Zjwknr9FEH5v/axRPdZwNafY4ltRQy1KaetPRwBF0Bg1EweGtUpjjh87lErJVvMonhWKvyICTxakVi136nwjA==	2024-05-18 11:25:04.469073	2024-05-18 11:25:04.469073	f	browser	{"nodes": [{"meta": {"label": {"id": 1010002, "text": "Sign in with discord", "type": "info", "context": {"provider": "discord"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "discord", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with facebook", "type": "info", "context": {"provider": "facebook"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "facebook", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with github", "type": "info", "context": {"provider": "github"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "github", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with google", "type": "info", "context": {"provider": "google"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "google", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with slack", "type": "info", "context": {"provider": "slack"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "slack", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with spotify", "type": "info", "context": {"provider": "spotify"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "spotify", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with twitch", "type": "info", "context": {"provider": "twitch"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "twitch", "disabled": false, "node_type": "input"}}, {"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "7Zjwknr9FEH5v/axRPdZwNafY4ltRQy1KaetPRwBF0Bg1EweGtUpjjh87lErJVvMonhWKvyICTxakVi136nwjA==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "E-Mail", "type": "info", "context": {"title": "E-Mail"}}}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "identifier", "type": "text", "value": "", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070001, "text": "Password", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "password", "type": "password", "disabled": false, "required": true, "node_type": "input", "autocomplete": "current-password"}}, {"meta": {"label": {"id": 1010001, "text": "Sign in", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "password", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/login?flow=753c9dc8-a04a-4655-8f9e-66cff204dc50", "method": "POST"}	e796004f-1c18-4960-9388-dcd5835741db	aal1	{}	\N	\N	choose_method	0	\N
6708f741-71ed-4ed0-b819-3cde079ff3bc	https://juicer-dev.xyz/self-service/login/browser	2024-05-18 11:25:59.77086	2024-05-18 11:35:59.77086		RZ8n0vTGpvdPcNXJO7OIeXzmqFzH9wQfEetTIG2NwYHgc7eKJsoqFrSCOrj7hPFFWpPzEEIPKUP3o1l6/Seixg==	2024-05-18 11:25:59.772348	2024-05-18 11:25:59.772348	f	browser	{"nodes": [{"meta": {"label": {"id": 1010002, "text": "Sign in with discord", "type": "info", "context": {"provider": "discord"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "discord", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with facebook", "type": "info", "context": {"provider": "facebook"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "facebook", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with github", "type": "info", "context": {"provider": "github"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "github", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with google", "type": "info", "context": {"provider": "google"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "google", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with slack", "type": "info", "context": {"provider": "slack"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "slack", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with spotify", "type": "info", "context": {"provider": "spotify"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "spotify", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with twitch", "type": "info", "context": {"provider": "twitch"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "twitch", "disabled": false, "node_type": "input"}}, {"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "RZ8n0vTGpvdPcNXJO7OIeXzmqFzH9wQfEetTIG2NwYHgc7eKJsoqFrSCOrj7hPFFWpPzEEIPKUP3o1l6/Seixg==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "E-Mail", "type": "info", "context": {"title": "E-Mail"}}}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "identifier", "type": "text", "value": "", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070001, "text": "Password", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "password", "type": "password", "disabled": false, "required": true, "node_type": "input", "autocomplete": "current-password"}}, {"meta": {"label": {"id": 1010001, "text": "Sign in", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "password", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/login?flow=6708f741-71ed-4ed0-b819-3cde079ff3bc", "method": "POST"}	e796004f-1c18-4960-9388-dcd5835741db	aal1	{}	\N	\N	choose_method	0	\N
488b918a-afc5-4a7f-88e0-1f7c8dcde850	https://juicer-dev.xyz/self-service/login/browser	2024-05-18 11:26:05.734366	2024-05-18 11:36:05.734366	password	V4e2cu2+MlhNpxGHyz/YQk9PxTNygqRy1SowzCJL6KXyayYqP7K+ubZV/vYLCKF+aTqef/d6iS4zYjqWsuGL4g==	2024-05-18 11:26:05.737891	2024-05-18 11:26:05.737891	f	browser	{"nodes": [{"meta": {"label": {"id": 1010002, "text": "Sign in with discord", "type": "info", "context": {"provider": "discord"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "discord", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with facebook", "type": "info", "context": {"provider": "facebook"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "facebook", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with github", "type": "info", "context": {"provider": "github"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "github", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with google", "type": "info", "context": {"provider": "google"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "google", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with slack", "type": "info", "context": {"provider": "slack"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "slack", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with spotify", "type": "info", "context": {"provider": "spotify"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "spotify", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1010002, "text": "Sign in with twitch", "type": "info", "context": {"provider": "twitch"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "twitch", "disabled": false, "node_type": "input"}}, {"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "V4e2cu2+MlhNpxGHyz/YQk9PxTNygqRy1SowzCJL6KXyayYqP7K+ubZV/vYLCKF+aTqef/d6iS4zYjqWsuGL4g==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "E-Mail", "type": "info", "context": {"title": "E-Mail"}}}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "identifier", "type": "text", "value": "", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070001, "text": "Password", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "password", "type": "password", "disabled": false, "required": true, "node_type": "input", "autocomplete": "current-password"}}, {"meta": {"label": {"id": 1010001, "text": "Sign in", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "password", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/login?flow=488b918a-afc5-4a7f-88e0-1f7c8dcde850", "method": "POST"}	e796004f-1c18-4960-9388-dcd5835741db	aal1	{}	\N	\N	choose_method	0	\N
\.


--
-- Data for Name: selfservice_recovery_flows; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.selfservice_recovery_flows (id, request_url, issued_at, expires_at, active_method, csrf_token, state, recovered_identity_id, created_at, updated_at, type, ui, nid, submit_count, skip_csrf_check) FROM stdin;
bf09d528-8833-43f9-a5ce-6c5c79145e1d	https://juicer-dev.xyz/self-service/recovery/browser	2024-05-18 11:24:33.541307	2024-05-18 12:24:33.541307	code	Vej4R3iQ15mD79WLI8PEQnOSH5Y/wHBH+ms1zSP7qs7YpETLGLjqVkIszWtMEcZOB3UqNa4Ndc6JXcBF4FNNAg==	choose_method	\N	2024-05-18 11:24:33.541421	2024-05-18 11:24:33.541421	browser	{"nodes": [{"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "Vej4R3iQ15mD79WLI8PEQnOSH5Y/wHBH+ms1zSP7qs7YpETLGLjqVkIszWtMEcZOB3UqNa4Ndc6JXcBF4FNNAg==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070007, "text": "Email", "type": "info"}}, "type": "input", "group": "code", "messages": [], "attributes": {"name": "email", "type": "email", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070005, "text": "Submit", "type": "info"}}, "type": "input", "group": "code", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "code", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/recovery?flow=bf09d528-8833-43f9-a5ce-6c5c79145e1d", "method": "POST"}	e796004f-1c18-4960-9388-dcd5835741db	0	f
\.


--
-- Data for Name: selfservice_registration_flows; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.selfservice_registration_flows (id, request_url, issued_at, expires_at, active_method, csrf_token, created_at, updated_at, type, ui, nid, internal_context, oauth2_login_challenge, oauth2_login_challenge_data, state, submit_count, organization_id) FROM stdin;
fdfe7f91-a68f-42fb-a66b-2d0471520c7c	https://juicer-dev.xyz/self-service/registration/browser	2024-05-18 11:24:30.86869	2024-05-18 11:34:30.86869		RBib4eGakz5xkXvBM1aSh2rQHmZx4ark6yGnf3itF4rJVCdtgbKu8bBSYyFchJCLHjcrxeAsr22YF1L3uwXwRg==	2024-05-18 11:24:30.870018	2024-05-18 11:24:30.870018	browser	{"nodes": [{"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "RBib4eGakz5xkXvBM1aSh2rQHmZx4ark6yGnf3itF4rJVCdtgbKu8bBSYyFchJCLHjcrxeAsr22YF1L3uwXwRg==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with discord", "type": "info", "context": {"provider": "discord"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "discord", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with facebook", "type": "info", "context": {"provider": "facebook"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "facebook", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with github", "type": "info", "context": {"provider": "github"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "github", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with google", "type": "info", "context": {"provider": "google"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "google", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with slack", "type": "info", "context": {"provider": "slack"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "slack", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with spotify", "type": "info", "context": {"provider": "spotify"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "spotify", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with twitch", "type": "info", "context": {"provider": "twitch"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "twitch", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "E-Mail", "type": "info", "context": {"title": "E-Mail"}}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.email", "type": "email", "disabled": false, "required": true, "node_type": "input", "autocomplete": "email"}}, {"meta": {"label": {"id": 1070001, "text": "Password", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "password", "type": "password", "disabled": false, "required": true, "node_type": "input", "autocomplete": "new-password"}}, {"meta": {"label": {"id": 1070002, "text": "First name", "type": "info", "context": {"title": "First name"}}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.first_name", "type": "text", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "Last name", "type": "info", "context": {"title": "Last name"}}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.last_name", "type": "text", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "Picture", "type": "info", "context": {"title": "Picture"}}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.avatar_url", "type": "text", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040001, "text": "Sign up", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "password", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/registration?flow=fdfe7f91-a68f-42fb-a66b-2d0471520c7c", "method": "POST"}	e796004f-1c18-4960-9388-dcd5835741db	{}	\N	\N	choose_method	0	\N
619c6561-50e5-4f6d-bdb7-4e56ec6cbfe8	https://juicer-dev.xyz/self-service/registration/browser	2024-05-18 11:25:22.609103	2024-05-18 11:35:22.609103		lCUyFlzAcZ44Bnz8X7AW1P9T7tgKj4u11v6IeR+CXuYZaY6aPOhMUfnFZBwwYhTYi7Tbe5tCjjylyH3x3Cq5Kg==	2024-05-18 11:25:22.61169	2024-05-18 11:25:22.61169	browser	{"nodes": [{"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "lCUyFlzAcZ44Bnz8X7AW1P9T7tgKj4u11v6IeR+CXuYZaY6aPOhMUfnFZBwwYhTYi7Tbe5tCjjylyH3x3Cq5Kg==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with discord", "type": "info", "context": {"provider": "discord"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "discord", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with facebook", "type": "info", "context": {"provider": "facebook"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "facebook", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with github", "type": "info", "context": {"provider": "github"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "github", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with google", "type": "info", "context": {"provider": "google"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "google", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with slack", "type": "info", "context": {"provider": "slack"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "slack", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with spotify", "type": "info", "context": {"provider": "spotify"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "spotify", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with twitch", "type": "info", "context": {"provider": "twitch"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "twitch", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "E-Mail", "type": "info", "context": {"title": "E-Mail"}}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.email", "type": "email", "disabled": false, "required": true, "node_type": "input", "autocomplete": "email"}}, {"meta": {"label": {"id": 1070001, "text": "Password", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "password", "type": "password", "disabled": false, "required": true, "node_type": "input", "autocomplete": "new-password"}}, {"meta": {"label": {"id": 1070002, "text": "First name", "type": "info", "context": {"title": "First name"}}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.first_name", "type": "text", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "Last name", "type": "info", "context": {"title": "Last name"}}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.last_name", "type": "text", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "Picture", "type": "info", "context": {"title": "Picture"}}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.avatar_url", "type": "text", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040001, "text": "Sign up", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "password", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/registration?flow=619c6561-50e5-4f6d-bdb7-4e56ec6cbfe8", "method": "POST"}	e796004f-1c18-4960-9388-dcd5835741db	{}	\N	\N	choose_method	0	\N
84bf4142-7cd9-4a65-b14c-e88d38d70d64	https://juicer-dev.xyz/self-service/registration/browser	2024-05-18 11:25:45.710556	2024-05-18 11:35:45.710556		aqiSAOldoJ3x430h9vvCDzYxeJEVmxY7SAFE4IuFsNLn5C6MiXWdUjAgZcGZKcADQtZNMoRWE7I7N7FoSC1XHg==	2024-05-18 11:25:45.715301	2024-05-18 11:25:45.715301	browser	{"nodes": [{"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "aqiSAOldoJ3x430h9vvCDzYxeJEVmxY7SAFE4IuFsNLn5C6MiXWdUjAgZcGZKcADQtZNMoRWE7I7N7FoSC1XHg==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with discord", "type": "info", "context": {"provider": "discord"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "discord", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with facebook", "type": "info", "context": {"provider": "facebook"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "facebook", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with github", "type": "info", "context": {"provider": "github"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "github", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with google", "type": "info", "context": {"provider": "google"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "google", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with slack", "type": "info", "context": {"provider": "slack"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "slack", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with spotify", "type": "info", "context": {"provider": "spotify"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "spotify", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040002, "text": "Sign up with twitch", "type": "info", "context": {"provider": "twitch"}}}, "type": "input", "group": "oidc", "messages": [], "attributes": {"name": "provider", "type": "submit", "value": "twitch", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "E-Mail", "type": "info", "context": {"title": "E-Mail"}}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.email", "type": "email", "disabled": false, "required": true, "node_type": "input", "autocomplete": "email"}}, {"meta": {"label": {"id": 1070001, "text": "Password", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "password", "type": "password", "disabled": false, "required": true, "node_type": "input", "autocomplete": "new-password"}}, {"meta": {"label": {"id": 1070002, "text": "First name", "type": "info", "context": {"title": "First name"}}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.first_name", "type": "text", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "Last name", "type": "info", "context": {"title": "Last name"}}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.last_name", "type": "text", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1070002, "text": "Picture", "type": "info", "context": {"title": "Picture"}}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "traits.avatar_url", "type": "text", "disabled": false, "node_type": "input"}}, {"meta": {"label": {"id": 1040001, "text": "Sign up", "type": "info"}}, "type": "input", "group": "password", "messages": [], "attributes": {"name": "method", "type": "submit", "value": "password", "disabled": false, "node_type": "input"}}], "action": "https://juicer-dev.xyz/kratos/self-service/registration?flow=84bf4142-7cd9-4a65-b14c-e88d38d70d64", "method": "POST"}	e796004f-1c18-4960-9388-dcd5835741db	{}	\N	\N	choose_method	0	\N
\.


--
-- Data for Name: selfservice_settings_flows; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.selfservice_settings_flows (id, request_url, issued_at, expires_at, identity_id, created_at, updated_at, active_method, state, type, ui, nid, internal_context) FROM stdin;
\.


--
-- Data for Name: selfservice_verification_flows; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.selfservice_verification_flows (id, request_url, issued_at, expires_at, csrf_token, created_at, updated_at, type, state, active_method, ui, nid, submit_count, oauth2_login_challenge, session_id, identity_id, authentication_methods) FROM stdin;
f38bd08b-a3ee-475b-8a24-eb89350a3b3f	https://juicer-dev.xyz/self-service/registration/browser	2024-05-18 11:25:45.991181	2024-05-18 12:25:45.991181	wvWljLXuhTT9Ana//Vgq0wlVB8dmDafvOgQW1USrkRdnGTXUZ+IJ1Qbwmc49b1PvLyBci+P1irPcTByP1AHyUA==	2024-05-18 11:25:45.991294	2024-05-18 11:25:45.991294	browser	passed_challenge	code	{"nodes": [{"meta": {}, "type": "input", "group": "default", "messages": [], "attributes": {"name": "csrf_token", "type": "hidden", "value": "wvWljLXuhTT9Ana//Vgq0wlVB8dmDafvOgQW1USrkRdnGTXUZ+IJ1Qbwmc49b1PvLyBci+P1irPcTByP1AHyUA==", "disabled": false, "required": true, "node_type": "input"}}, {"meta": {"label": {"id": 1070009, "text": "Continue", "type": "info"}}, "type": "a", "group": "code", "messages": [], "attributes": {"id": "continue", "href": "https://juicer-dev.xyz/auth/login", "title": {"id": 1070009, "text": "Continue", "type": "info"}, "node_type": "a"}}], "action": "https://juicer-dev.xyz/auth/login", "method": "GET", "messages": [{"id": 1080002, "text": "You successfully verified your email address.", "type": "success"}]}	e796004f-1c18-4960-9388-dcd5835741db	1	\N	89699db2-5162-4ddf-8ae0-9c1cb369a560	bb415d23-d40a-405d-8279-10fab2c6d9a0	[{"method":"password","aal":"aal1","completed_at":"2024-05-18T11:25:45.985282398Z"}]
\.


--
-- Data for Name: session_devices; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.session_devices (id, ip_address, user_agent, location, nid, session_id, created_at, updated_at) FROM stdin;
647c511c-5855-4629-9957-6195d6eeb471	172.19.0.1	Mozilla/5.0 (X11; Linux x86_64; rv:125.0) Gecko/20100101 Firefox/125.0		e796004f-1c18-4960-9388-dcd5835741db	89699db2-5162-4ddf-8ae0-9c1cb369a560	2024-05-18 11:25:45.987481	2024-05-18 11:25:45.987481
e6b0db9d-9f65-4fcb-bc7b-239d69cbb003	172.19.0.1	Mozilla/5.0 (X11; Linux x86_64; rv:125.0) Gecko/20100101 Firefox/125.0		e796004f-1c18-4960-9388-dcd5835741db	ea4c3a90-0de5-4a1b-80bd-720c25d788a6	2024-05-18 11:26:05.997376	2024-05-18 11:26:05.997376
\.


--
-- Data for Name: session_token_exchanges; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.session_token_exchanges (id, nid, flow_id, session_id, init_code, return_to_code, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.sessions (id, issued_at, expires_at, authenticated_at, identity_id, created_at, updated_at, token, active, nid, logout_token, aal, authentication_methods) FROM stdin;
89699db2-5162-4ddf-8ae0-9c1cb369a560	2024-05-18 11:25:45.985282	2024-05-19 11:25:45.985282	2024-05-18 11:25:45.98529	bb415d23-d40a-405d-8279-10fab2c6d9a0	2024-05-18 11:25:45.985461	2024-05-18 11:25:45.985461	ory_st_vr6FqKptPopNJWzGMPVwB95GedZsG2MI	t	e796004f-1c18-4960-9388-dcd5835741db	ory_lo_QNyLXwMFTDntUy3ce1SYCTMBTIyQLbJI	aal1	[{"aal": "aal1", "method": "password", "completed_at": "2024-05-18T11:25:45.985282398Z"}]
ea4c3a90-0de5-4a1b-80bd-720c25d788a6	2024-05-18 11:26:05.994303	2024-05-19 11:26:05.994303	2024-05-18 11:26:05.994303	bb415d23-d40a-405d-8279-10fab2c6d9a0	2024-05-18 11:26:05.994711	2024-05-18 11:26:05.994711	ory_st_N5v1cBTcwy9tAM2dU5ugQDpdcu0CTc9P	t	e796004f-1c18-4960-9388-dcd5835741db	ory_lo_Tig4958FFQtNxYIMbHFNcbAtxt6jJhuK	aal1	[{"aal": "aal1", "method": "password", "completed_at": "2024-05-18T11:26:05.99428343Z"}]
\.


--
-- Name: continuity_containers continuity_containers_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.continuity_containers
    ADD CONSTRAINT continuity_containers_pkey PRIMARY KEY (id);


--
-- Name: courier_message_dispatches courier_message_dispatches_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.courier_message_dispatches
    ADD CONSTRAINT courier_message_dispatches_pkey PRIMARY KEY (id);


--
-- Name: courier_messages courier_messages_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.courier_messages
    ADD CONSTRAINT courier_messages_pkey PRIMARY KEY (id);


--
-- Name: identities identities_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identities
    ADD CONSTRAINT identities_pkey PRIMARY KEY (id);


--
-- Name: identity_credential_identifiers identity_credential_identifiers_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_credential_identifiers
    ADD CONSTRAINT identity_credential_identifiers_pkey PRIMARY KEY (id);


--
-- Name: identity_credential_types identity_credential_types_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_credential_types
    ADD CONSTRAINT identity_credential_types_pkey PRIMARY KEY (id);


--
-- Name: identity_credentials identity_credentials_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_credentials
    ADD CONSTRAINT identity_credentials_pkey PRIMARY KEY (id);


--
-- Name: identity_login_codes identity_login_codes_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_login_codes
    ADD CONSTRAINT identity_login_codes_pkey PRIMARY KEY (id);


--
-- Name: identity_recovery_addresses identity_recovery_addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_addresses
    ADD CONSTRAINT identity_recovery_addresses_pkey PRIMARY KEY (id);


--
-- Name: identity_recovery_codes identity_recovery_codes_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_codes
    ADD CONSTRAINT identity_recovery_codes_pkey PRIMARY KEY (id);


--
-- Name: identity_recovery_tokens identity_recovery_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_tokens
    ADD CONSTRAINT identity_recovery_tokens_pkey PRIMARY KEY (id);


--
-- Name: identity_registration_codes identity_registration_codes_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_registration_codes
    ADD CONSTRAINT identity_registration_codes_pkey PRIMARY KEY (id);


--
-- Name: identity_verifiable_addresses identity_verifiable_addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verifiable_addresses
    ADD CONSTRAINT identity_verifiable_addresses_pkey PRIMARY KEY (id);


--
-- Name: identity_verification_codes identity_verification_codes_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verification_codes
    ADD CONSTRAINT identity_verification_codes_pkey PRIMARY KEY (id);


--
-- Name: identity_verification_tokens identity_verification_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verification_tokens
    ADD CONSTRAINT identity_verification_tokens_pkey PRIMARY KEY (id);


--
-- Name: keto_relation_tuples keto_relation_tuples_uuid_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.keto_relation_tuples
    ADD CONSTRAINT keto_relation_tuples_uuid_pkey PRIMARY KEY (shard_id, nid);


--
-- Name: keto_uuid_mappings keto_uuid_mappings_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.keto_uuid_mappings
    ADD CONSTRAINT keto_uuid_mappings_pkey PRIMARY KEY (id);


--
-- Name: networks networks_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.networks
    ADD CONSTRAINT networks_pkey PRIMARY KEY (id);


--
-- Name: selfservice_errors selfservice_errors_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_errors
    ADD CONSTRAINT selfservice_errors_pkey PRIMARY KEY (id);


--
-- Name: selfservice_login_flows selfservice_login_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_login_flows
    ADD CONSTRAINT selfservice_login_requests_pkey PRIMARY KEY (id);


--
-- Name: selfservice_settings_flows selfservice_profile_management_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_settings_flows
    ADD CONSTRAINT selfservice_profile_management_requests_pkey PRIMARY KEY (id);


--
-- Name: selfservice_recovery_flows selfservice_recovery_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_recovery_flows
    ADD CONSTRAINT selfservice_recovery_requests_pkey PRIMARY KEY (id);


--
-- Name: selfservice_registration_flows selfservice_registration_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_registration_flows
    ADD CONSTRAINT selfservice_registration_requests_pkey PRIMARY KEY (id);


--
-- Name: selfservice_verification_flows selfservice_verification_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_verification_flows
    ADD CONSTRAINT selfservice_verification_requests_pkey PRIMARY KEY (id);


--
-- Name: session_devices session_devices_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.session_devices
    ADD CONSTRAINT session_devices_pkey PRIMARY KEY (id);


--
-- Name: session_token_exchanges session_token_exchanges_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.session_token_exchanges
    ADD CONSTRAINT session_token_exchanges_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- Name: session_devices unique_session_device; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.session_devices
    ADD CONSTRAINT unique_session_device UNIQUE (nid, session_id, ip_address, user_agent);


--
-- Name: continuity_containers_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX continuity_containers_id_nid_idx ON public.continuity_containers USING btree (id, nid);


--
-- Name: continuity_containers_identity_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX continuity_containers_identity_id_nid_idx ON public.continuity_containers USING btree (identity_id, nid);


--
-- Name: continuity_containers_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX continuity_containers_nid_id_idx ON public.continuity_containers USING btree (nid, id);


--
-- Name: courier_message_dispatches_id_message_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX courier_message_dispatches_id_message_id_nid_idx ON public.courier_message_dispatches USING btree (id, message_id, nid);


--
-- Name: courier_messages_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX courier_messages_id_nid_idx ON public.courier_messages USING btree (id, nid);


--
-- Name: courier_messages_nid_created_at_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX courier_messages_nid_created_at_id_idx ON public.courier_messages USING btree (nid, created_at DESC);


--
-- Name: courier_messages_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX courier_messages_nid_id_idx ON public.courier_messages USING btree (nid, id);


--
-- Name: courier_messages_nid_recipient_created_at_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX courier_messages_nid_recipient_created_at_id_idx ON public.courier_messages USING btree (nid, recipient, created_at DESC);


--
-- Name: courier_messages_nid_status_created_at_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX courier_messages_nid_status_created_at_id_idx ON public.courier_messages USING btree (nid, status, created_at DESC);


--
-- Name: courier_messages_status_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX courier_messages_status_idx ON public.courier_messages USING btree (status);


--
-- Name: identities_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identities_id_nid_idx ON public.identities USING btree (id, nid);


--
-- Name: identities_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identities_nid_id_idx ON public.identities USING btree (nid, id);


--
-- Name: identity_credential_identifiers_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_credential_identifiers_id_nid_idx ON public.identity_credential_identifiers USING btree (id, nid);


--
-- Name: identity_credential_identifiers_identifier_nid_type_uq_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX identity_credential_identifiers_identifier_nid_type_uq_idx ON public.identity_credential_identifiers USING btree (nid, identity_credential_type_id, identifier);


--
-- Name: identity_credential_identifiers_nid_i_ici_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_credential_identifiers_nid_i_ici_idx ON public.identity_credential_identifiers USING btree (nid, identifier, identity_credential_id);


--
-- Name: identity_credential_identifiers_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_credential_identifiers_nid_id_idx ON public.identity_credential_identifiers USING btree (nid, id);


--
-- Name: identity_credential_identifiers_nid_identifier_gin; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_credential_identifiers_nid_identifier_gin ON public.identity_credential_identifiers USING gin (nid, identifier public.gin_trgm_ops);


--
-- Name: identity_credential_identifiers_nid_identity_credential_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_credential_identifiers_nid_identity_credential_id_idx ON public.identity_credential_identifiers USING btree (identity_credential_id, nid);


--
-- Name: identity_credential_types_name_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX identity_credential_types_name_idx ON public.identity_credential_types USING btree (name);


--
-- Name: identity_credentials_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_credentials_id_nid_idx ON public.identity_credentials USING btree (id, nid);


--
-- Name: identity_credentials_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_credentials_nid_id_idx ON public.identity_credentials USING btree (nid, id);


--
-- Name: identity_credentials_nid_identity_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_credentials_nid_identity_id_idx ON public.identity_credentials USING btree (identity_id, nid);


--
-- Name: identity_login_codes_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_login_codes_id_nid_idx ON public.identity_login_codes USING btree (id, nid);


--
-- Name: identity_login_codes_nid_flow_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_login_codes_nid_flow_id_idx ON public.identity_login_codes USING btree (nid, selfservice_login_flow_id);


--
-- Name: identity_recovery_addresses_code_uq_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX identity_recovery_addresses_code_uq_idx ON public.identity_recovery_tokens USING btree (token);


--
-- Name: identity_recovery_addresses_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_addresses_id_nid_idx ON public.identity_recovery_addresses USING btree (id, nid);


--
-- Name: identity_recovery_addresses_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_addresses_nid_id_idx ON public.identity_recovery_addresses USING btree (nid, id);


--
-- Name: identity_recovery_addresses_nid_identity_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_addresses_nid_identity_id_idx ON public.identity_recovery_addresses USING btree (identity_id, nid);


--
-- Name: identity_recovery_addresses_status_via_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_addresses_status_via_idx ON public.identity_recovery_addresses USING btree (nid, via, value);


--
-- Name: identity_recovery_addresses_status_via_uq_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX identity_recovery_addresses_status_via_uq_idx ON public.identity_recovery_addresses USING btree (nid, via, value);


--
-- Name: identity_recovery_codes_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_codes_id_nid_idx ON public.identity_recovery_codes USING btree (id, nid);


--
-- Name: identity_recovery_codes_identity_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_codes_identity_id_nid_idx ON public.identity_recovery_codes USING btree (identity_id, nid);


--
-- Name: identity_recovery_codes_identity_recovery_address_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_codes_identity_recovery_address_id_nid_idx ON public.identity_recovery_codes USING btree (identity_recovery_address_id, nid);


--
-- Name: identity_recovery_codes_nid_flow_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_codes_nid_flow_id_idx ON public.identity_recovery_codes USING btree (nid, selfservice_recovery_flow_id);


--
-- Name: identity_recovery_tokens_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_tokens_id_nid_idx ON public.identity_recovery_tokens USING btree (id, nid);


--
-- Name: identity_recovery_tokens_identity_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_tokens_identity_id_nid_idx ON public.identity_recovery_tokens USING btree (identity_id, nid);


--
-- Name: identity_recovery_tokens_identity_recovery_address_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_tokens_identity_recovery_address_id_idx ON public.identity_recovery_tokens USING btree (identity_recovery_address_id);


--
-- Name: identity_recovery_tokens_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_tokens_nid_id_idx ON public.identity_recovery_tokens USING btree (nid, id);


--
-- Name: identity_recovery_tokens_selfservice_recovery_flow_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_tokens_selfservice_recovery_flow_id_idx ON public.identity_recovery_tokens USING btree (selfservice_recovery_flow_id);


--
-- Name: identity_recovery_tokens_token_nid_used_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_recovery_tokens_token_nid_used_idx ON public.identity_recovery_tokens USING btree (nid, token, used);


--
-- Name: identity_registration_codes_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_registration_codes_id_nid_idx ON public.identity_registration_codes USING btree (id, nid);


--
-- Name: identity_registration_codes_nid_flow_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_registration_codes_nid_flow_id_idx ON public.identity_registration_codes USING btree (nid, selfservice_registration_flow_id);


--
-- Name: identity_verifiable_addresses_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verifiable_addresses_id_nid_idx ON public.identity_verifiable_addresses USING btree (id, nid);


--
-- Name: identity_verifiable_addresses_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verifiable_addresses_nid_id_idx ON public.identity_verifiable_addresses USING btree (nid, id);


--
-- Name: identity_verifiable_addresses_nid_identity_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verifiable_addresses_nid_identity_id_idx ON public.identity_verifiable_addresses USING btree (identity_id, nid);


--
-- Name: identity_verifiable_addresses_status_via_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verifiable_addresses_status_via_idx ON public.identity_verifiable_addresses USING btree (nid, via, value);


--
-- Name: identity_verifiable_addresses_status_via_uq_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX identity_verifiable_addresses_status_via_uq_idx ON public.identity_verifiable_addresses USING btree (nid, via, value);


--
-- Name: identity_verification_codes_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verification_codes_id_nid_idx ON public.identity_verification_codes USING btree (id, nid);


--
-- Name: identity_verification_codes_nid_flow_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verification_codes_nid_flow_id_idx ON public.identity_verification_codes USING btree (nid, selfservice_verification_flow_id);


--
-- Name: identity_verification_codes_verifiable_address_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verification_codes_verifiable_address_nid_idx ON public.identity_verification_codes USING btree (identity_verifiable_address_id, nid);


--
-- Name: identity_verification_tokens_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verification_tokens_id_nid_idx ON public.identity_verification_tokens USING btree (id, nid);


--
-- Name: identity_verification_tokens_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verification_tokens_nid_id_idx ON public.identity_verification_tokens USING btree (nid, id);


--
-- Name: identity_verification_tokens_token_nid_used_flow_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verification_tokens_token_nid_used_flow_id_idx ON public.identity_verification_tokens USING btree (nid, token, used, selfservice_verification_flow_id);


--
-- Name: identity_verification_tokens_token_uq_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX identity_verification_tokens_token_uq_idx ON public.identity_verification_tokens USING btree (token);


--
-- Name: identity_verification_tokens_verifiable_address_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verification_tokens_verifiable_address_id_idx ON public.identity_verification_tokens USING btree (identity_verifiable_address_id);


--
-- Name: identity_verification_tokens_verification_flow_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX identity_verification_tokens_verification_flow_id_idx ON public.identity_verification_tokens USING btree (selfservice_verification_flow_id);


--
-- Name: keto_relation_tuples_uuid_full_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX keto_relation_tuples_uuid_full_idx ON public.keto_relation_tuples USING btree (nid, namespace, object, relation, subject_id, subject_set_namespace, subject_set_object, subject_set_relation, commit_time);


--
-- Name: keto_relation_tuples_uuid_reverse_subject_ids_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX keto_relation_tuples_uuid_reverse_subject_ids_idx ON public.keto_relation_tuples USING btree (nid, subject_id, relation, namespace) WHERE ((subject_set_namespace IS NULL) AND (subject_set_object IS NULL) AND (subject_set_relation IS NULL));


--
-- Name: keto_relation_tuples_uuid_reverse_subject_sets_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX keto_relation_tuples_uuid_reverse_subject_sets_idx ON public.keto_relation_tuples USING btree (nid, subject_set_namespace, subject_set_object, subject_set_relation, relation, namespace) WHERE (subject_id IS NULL);


--
-- Name: keto_relation_tuples_uuid_subject_ids_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX keto_relation_tuples_uuid_subject_ids_idx ON public.keto_relation_tuples USING btree (nid, namespace, object, relation, subject_id) WHERE ((subject_set_namespace IS NULL) AND (subject_set_object IS NULL) AND (subject_set_relation IS NULL));


--
-- Name: keto_relation_tuples_uuid_subject_sets_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX keto_relation_tuples_uuid_subject_sets_idx ON public.keto_relation_tuples USING btree (nid, namespace, object, relation, subject_set_namespace, subject_set_object, subject_set_relation) WHERE (subject_id IS NULL);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: schema_migration_version_self_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX schema_migration_version_self_idx ON public.schema_migration USING btree (version_self);


--
-- Name: selfservice_errors_errors_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_errors_errors_nid_id_idx ON public.selfservice_errors USING btree (nid, id);


--
-- Name: selfservice_login_flows_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_login_flows_id_nid_idx ON public.selfservice_login_flows USING btree (id, nid);


--
-- Name: selfservice_login_flows_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_login_flows_nid_id_idx ON public.selfservice_login_flows USING btree (nid, id);


--
-- Name: selfservice_recovery_flows_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_recovery_flows_id_nid_idx ON public.selfservice_recovery_flows USING btree (id, nid);


--
-- Name: selfservice_recovery_flows_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_recovery_flows_nid_id_idx ON public.selfservice_recovery_flows USING btree (nid, id);


--
-- Name: selfservice_recovery_flows_recovered_identity_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_recovery_flows_recovered_identity_id_nid_idx ON public.selfservice_recovery_flows USING btree (recovered_identity_id, nid);


--
-- Name: selfservice_registration_flows_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_registration_flows_id_nid_idx ON public.selfservice_registration_flows USING btree (id, nid);


--
-- Name: selfservice_registration_flows_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_registration_flows_nid_id_idx ON public.selfservice_registration_flows USING btree (nid, id);


--
-- Name: selfservice_settings_flows_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_settings_flows_id_nid_idx ON public.selfservice_settings_flows USING btree (id, nid);


--
-- Name: selfservice_settings_flows_identity_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_settings_flows_identity_id_nid_idx ON public.selfservice_settings_flows USING btree (identity_id, nid);


--
-- Name: selfservice_settings_flows_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_settings_flows_nid_id_idx ON public.selfservice_settings_flows USING btree (nid, id);


--
-- Name: selfservice_verification_flows_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_verification_flows_id_nid_idx ON public.selfservice_verification_flows USING btree (id, nid);


--
-- Name: selfservice_verification_flows_nid_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX selfservice_verification_flows_nid_id_idx ON public.selfservice_verification_flows USING btree (nid, id);


--
-- Name: session_devices_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX session_devices_id_nid_idx ON public.session_devices USING btree (id, nid);


--
-- Name: session_devices_session_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX session_devices_session_id_nid_idx ON public.session_devices USING btree (session_id, nid);


--
-- Name: session_token_exchanges_nid_code_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX session_token_exchanges_nid_code_idx ON public.session_token_exchanges USING btree (init_code, nid);


--
-- Name: session_token_exchanges_nid_flow_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX session_token_exchanges_nid_flow_id_idx ON public.session_token_exchanges USING btree (flow_id, nid);


--
-- Name: sessions_id_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX sessions_id_nid_idx ON public.sessions USING btree (id, nid);


--
-- Name: sessions_identity_id_nid_sorted_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX sessions_identity_id_nid_sorted_idx ON public.sessions USING btree (identity_id, nid, authenticated_at DESC);


--
-- Name: sessions_logout_token_uq_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX sessions_logout_token_uq_idx ON public.sessions USING btree (logout_token);


--
-- Name: sessions_nid_created_at_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX sessions_nid_created_at_id_idx ON public.sessions USING btree (nid, created_at DESC, id);


--
-- Name: sessions_nid_id_identity_id_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX sessions_nid_id_identity_id_idx ON public.sessions USING btree (nid, identity_id, id);


--
-- Name: sessions_token_nid_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX sessions_token_nid_idx ON public.sessions USING btree (nid, token);


--
-- Name: sessions_token_uq_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX sessions_token_uq_idx ON public.sessions USING btree (token);


--
-- Name: continuity_containers continuity_containers_identity_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.continuity_containers
    ADD CONSTRAINT continuity_containers_identity_id_fkey FOREIGN KEY (identity_id) REFERENCES public.identities(id) ON DELETE CASCADE;


--
-- Name: continuity_containers continuity_containers_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.continuity_containers
    ADD CONSTRAINT continuity_containers_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: courier_message_dispatches courier_message_dispatches_message_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.courier_message_dispatches
    ADD CONSTRAINT courier_message_dispatches_message_id_fk FOREIGN KEY (message_id) REFERENCES public.courier_messages(id) ON DELETE CASCADE;


--
-- Name: courier_message_dispatches courier_message_dispatches_nid_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.courier_message_dispatches
    ADD CONSTRAINT courier_message_dispatches_nid_fk FOREIGN KEY (nid) REFERENCES public.networks(id) ON DELETE CASCADE;


--
-- Name: courier_messages courier_messages_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.courier_messages
    ADD CONSTRAINT courier_messages_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identities identities_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identities
    ADD CONSTRAINT identities_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_credential_identifiers identity_credential_identifiers_identity_credential_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_credential_identifiers
    ADD CONSTRAINT identity_credential_identifiers_identity_credential_id_fkey FOREIGN KEY (identity_credential_id) REFERENCES public.identity_credentials(id) ON DELETE CASCADE;


--
-- Name: identity_credential_identifiers identity_credential_identifiers_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_credential_identifiers
    ADD CONSTRAINT identity_credential_identifiers_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_credential_identifiers identity_credential_identifiers_type_id_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_credential_identifiers
    ADD CONSTRAINT identity_credential_identifiers_type_id_fk_idx FOREIGN KEY (identity_credential_type_id) REFERENCES public.identity_credential_types(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_credentials identity_credentials_identity_credential_type_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_credentials
    ADD CONSTRAINT identity_credentials_identity_credential_type_id_fkey FOREIGN KEY (identity_credential_type_id) REFERENCES public.identity_credential_types(id) ON DELETE CASCADE;


--
-- Name: identity_credentials identity_credentials_identity_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_credentials
    ADD CONSTRAINT identity_credentials_identity_id_fkey FOREIGN KEY (identity_id) REFERENCES public.identities(id) ON DELETE CASCADE;


--
-- Name: identity_credentials identity_credentials_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_credentials
    ADD CONSTRAINT identity_credentials_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_login_codes identity_login_codes_networks_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_login_codes
    ADD CONSTRAINT identity_login_codes_networks_id_fk FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_login_codes identity_login_codes_selfservice_login_flows_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_login_codes
    ADD CONSTRAINT identity_login_codes_selfservice_login_flows_id_fk FOREIGN KEY (selfservice_login_flow_id) REFERENCES public.selfservice_login_flows(id) ON DELETE CASCADE;


--
-- Name: identity_recovery_addresses identity_recovery_addresses_identity_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_addresses
    ADD CONSTRAINT identity_recovery_addresses_identity_id_fkey FOREIGN KEY (identity_id) REFERENCES public.identities(id) ON DELETE CASCADE;


--
-- Name: identity_recovery_addresses identity_recovery_addresses_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_addresses
    ADD CONSTRAINT identity_recovery_addresses_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_recovery_codes identity_recovery_codes_identity_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_codes
    ADD CONSTRAINT identity_recovery_codes_identity_id_fk FOREIGN KEY (identity_id) REFERENCES public.identities(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_recovery_codes identity_recovery_codes_identity_recovery_addresses_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_codes
    ADD CONSTRAINT identity_recovery_codes_identity_recovery_addresses_id_fk FOREIGN KEY (identity_recovery_address_id) REFERENCES public.identity_recovery_addresses(id) ON DELETE CASCADE;


--
-- Name: identity_recovery_codes identity_recovery_codes_networks_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_codes
    ADD CONSTRAINT identity_recovery_codes_networks_id_fk FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_recovery_codes identity_recovery_codes_selfservice_recovery_flows_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_codes
    ADD CONSTRAINT identity_recovery_codes_selfservice_recovery_flows_id_fk FOREIGN KEY (selfservice_recovery_flow_id) REFERENCES public.selfservice_recovery_flows(id) ON DELETE CASCADE;


--
-- Name: identity_recovery_tokens identity_recovery_tokens_identity_id_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_tokens
    ADD CONSTRAINT identity_recovery_tokens_identity_id_fk_idx FOREIGN KEY (identity_id) REFERENCES public.identities(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_recovery_tokens identity_recovery_tokens_identity_recovery_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_tokens
    ADD CONSTRAINT identity_recovery_tokens_identity_recovery_address_id_fkey FOREIGN KEY (identity_recovery_address_id) REFERENCES public.identity_recovery_addresses(id) ON DELETE CASCADE;


--
-- Name: identity_recovery_tokens identity_recovery_tokens_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_tokens
    ADD CONSTRAINT identity_recovery_tokens_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_recovery_tokens identity_recovery_tokens_selfservice_recovery_request_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_recovery_tokens
    ADD CONSTRAINT identity_recovery_tokens_selfservice_recovery_request_id_fkey FOREIGN KEY (selfservice_recovery_flow_id) REFERENCES public.selfservice_recovery_flows(id) ON DELETE CASCADE;


--
-- Name: identity_registration_codes identity_registration_codes_networks_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_registration_codes
    ADD CONSTRAINT identity_registration_codes_networks_id_fk FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_registration_codes identity_registration_codes_selfservice_registration_flows_id_f; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_registration_codes
    ADD CONSTRAINT identity_registration_codes_selfservice_registration_flows_id_f FOREIGN KEY (selfservice_registration_flow_id) REFERENCES public.selfservice_registration_flows(id) ON DELETE CASCADE;


--
-- Name: identity_verifiable_addresses identity_verifiable_addresses_identity_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verifiable_addresses
    ADD CONSTRAINT identity_verifiable_addresses_identity_id_fkey FOREIGN KEY (identity_id) REFERENCES public.identities(id) ON DELETE CASCADE;


--
-- Name: identity_verifiable_addresses identity_verifiable_addresses_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verifiable_addresses
    ADD CONSTRAINT identity_verifiable_addresses_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_verification_codes identity_verification_codes_identity_verifiable_addresses_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verification_codes
    ADD CONSTRAINT identity_verification_codes_identity_verifiable_addresses_id_fk FOREIGN KEY (identity_verifiable_address_id) REFERENCES public.identity_verifiable_addresses(id) ON DELETE CASCADE;


--
-- Name: identity_verification_codes identity_verification_codes_networks_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verification_codes
    ADD CONSTRAINT identity_verification_codes_networks_id_fk FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_verification_codes identity_verification_codes_selfservice_verification_flows_id_f; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verification_codes
    ADD CONSTRAINT identity_verification_codes_selfservice_verification_flows_id_f FOREIGN KEY (selfservice_verification_flow_id) REFERENCES public.selfservice_verification_flows(id) ON DELETE CASCADE;


--
-- Name: identity_verification_tokens identity_verification_tokens_identity_verifiable_address_i_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verification_tokens
    ADD CONSTRAINT identity_verification_tokens_identity_verifiable_address_i_fkey FOREIGN KEY (identity_verifiable_address_id) REFERENCES public.identity_verifiable_addresses(id) ON DELETE CASCADE;


--
-- Name: identity_verification_tokens identity_verification_tokens_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verification_tokens
    ADD CONSTRAINT identity_verification_tokens_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: identity_verification_tokens identity_verification_tokens_selfservice_verification_flow_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.identity_verification_tokens
    ADD CONSTRAINT identity_verification_tokens_selfservice_verification_flow_fkey FOREIGN KEY (selfservice_verification_flow_id) REFERENCES public.selfservice_verification_flows(id) ON DELETE CASCADE;


--
-- Name: keto_relation_tuples keto_relation_tuples_nid_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.keto_relation_tuples
    ADD CONSTRAINT keto_relation_tuples_nid_fk FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: selfservice_errors selfservice_errors_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_errors
    ADD CONSTRAINT selfservice_errors_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: selfservice_login_flows selfservice_login_flows_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_login_flows
    ADD CONSTRAINT selfservice_login_flows_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: selfservice_settings_flows selfservice_profile_management_requests_identity_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_settings_flows
    ADD CONSTRAINT selfservice_profile_management_requests_identity_id_fkey FOREIGN KEY (identity_id) REFERENCES public.identities(id) ON DELETE CASCADE;


--
-- Name: selfservice_recovery_flows selfservice_recovery_flows_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_recovery_flows
    ADD CONSTRAINT selfservice_recovery_flows_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: selfservice_recovery_flows selfservice_recovery_requests_recovered_identity_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_recovery_flows
    ADD CONSTRAINT selfservice_recovery_requests_recovered_identity_id_fkey FOREIGN KEY (recovered_identity_id) REFERENCES public.identities(id) ON DELETE CASCADE;


--
-- Name: selfservice_registration_flows selfservice_registration_flows_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_registration_flows
    ADD CONSTRAINT selfservice_registration_flows_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: selfservice_settings_flows selfservice_settings_flows_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_settings_flows
    ADD CONSTRAINT selfservice_settings_flows_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: selfservice_verification_flows selfservice_verification_flows_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.selfservice_verification_flows
    ADD CONSTRAINT selfservice_verification_flows_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: session_devices session_metadata_nid_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.session_devices
    ADD CONSTRAINT session_metadata_nid_fk FOREIGN KEY (nid) REFERENCES public.networks(id) ON DELETE CASCADE;


--
-- Name: session_devices session_metadata_sessions_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.session_devices
    ADD CONSTRAINT session_metadata_sessions_id_fk FOREIGN KEY (session_id) REFERENCES public.sessions(id) ON DELETE CASCADE;


--
-- Name: sessions sessions_identity_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_identity_id_fkey FOREIGN KEY (identity_id) REFERENCES public.identities(id) ON DELETE CASCADE;


--
-- Name: sessions sessions_nid_fk_idx; Type: FK CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_nid_fk_idx FOREIGN KEY (nid) REFERENCES public.networks(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

--
-- PostgreSQL database cluster dump complete
--

