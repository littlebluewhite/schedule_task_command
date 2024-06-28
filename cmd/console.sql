UPDATE stage_item
SET mode = 'monitor'
WHERE command_template_id IN (110, 113, 111, 142);