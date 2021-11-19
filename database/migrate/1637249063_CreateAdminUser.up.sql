-- Migration goes here.

INSERT INTO public.users (id,username,user_full_name,email,cpf_cnpj,"password",role_id,created_at,updated_at)
VALUES (gen_random_uuid ()::uuid,'admin','Admin','admin@viniciussantos.dev','44611032850','$2a$14$ayByQwQ2fxaKvx1KZDBW0ezvBCKROdkiW9GLKDHTVp1S7nl0wh3ZO',3,'NOW()','NOW()');