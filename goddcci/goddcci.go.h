#include <ddccontrol/ddcci.h>

char* xmlCharToChar(xmlChar* c){
    return (char*)c;
}


static int find_write_delay(struct monitor* mon, char ctrl) {

	struct monitor_db* monitor = mon->db;
	struct group_db* group;
	struct subgroup_db* subgroup;
	struct control_db* control;

	if (monitor)
	{
		/* loop through groups */
		for (group = monitor->group_list; group != NULL; group = group->next)
		{
			/* loop through subgroups inside group */
			for (subgroup = group->subgroup_list; subgroup != NULL; subgroup = subgroup->next)
			{
				/* loop through controls inside subgroup */
				for (control = subgroup->control_list; control != NULL; control = control->next)
				{
					/* check for control id */
					if (control->address == ctrl)
					{
						return control->delay;
					}
				}
			}
		}
	}
	return -1;
}