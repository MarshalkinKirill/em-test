--
-- PostgreSQL database dump
--

-- Dumped from database version 15.4
-- Dumped by pg_dump version 15.4

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
-- Name: emsrv; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA emsrv;


ALTER SCHEMA emsrv OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: people; Type: TABLE; Schema: emsrv; Owner: postgres
--

CREATE TABLE emsrv.people (
    "personId" uuid NOT NULL,
    name character varying(25) NOT NULL,
    surname character varying(25) NOT NULL,
    patronymic character varying(25) NOT NULL,
    age integer,
    gender character varying(10),
    nationality character varying(10)
);


ALTER TABLE emsrv.people OWNER TO postgres;

ALTER TABLE ONLY emsrv.people
    ADD CONSTRAINT person_pkey PRIMARY KEY ("personId");

