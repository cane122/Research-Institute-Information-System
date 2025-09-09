export namespace models {
	
	export class Korisnici {
	    korisnik_id: number;
	    korisnicko_ime: string;
	    email: string;
	    ime?: string;
	    prezime?: string;
	    uloga_id: number;
	    status: string;
	    poslednja_prijava?: string;
	    kreiran_datuma: string;
	    naziv_uloge?: string;
	
	    static createFrom(source: any = {}) {
	        return new Korisnici(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.korisnik_id = source["korisnik_id"];
	        this.korisnicko_ime = source["korisnicko_ime"];
	        this.email = source["email"];
	        this.ime = source["ime"];
	        this.prezime = source["prezime"];
	        this.uloga_id = source["uloga_id"];
	        this.status = source["status"];
	        this.poslednja_prijava = source["poslednja_prijava"];
	        this.kreiran_datuma = source["kreiran_datuma"];
	        this.naziv_uloge = source["naziv_uloge"];
	    }
	}
	export class Projekti {
	    projekat_id: number;
	    naziv_projekta: string;
	    opis?: string;
	    datum_pocetka?: string;
	    datum_zavrsetka?: string;
	    status: string;
	    rukovodilac_id?: number;
	    radni_tok_id?: number;
	    rukovodilac_ime?: string;
	    broj_zadataka?: number;
	    broj_clanova?: number;
	
	    static createFrom(source: any = {}) {
	        return new Projekti(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.projekat_id = source["projekat_id"];
	        this.naziv_projekta = source["naziv_projekta"];
	        this.opis = source["opis"];
	        this.datum_pocetka = source["datum_pocetka"];
	        this.datum_zavrsetka = source["datum_zavrsetka"];
	        this.status = source["status"];
	        this.rukovodilac_id = source["rukovodilac_id"];
	        this.radni_tok_id = source["radni_tok_id"];
	        this.rukovodilac_ime = source["rukovodilac_ime"];
	        this.broj_zadataka = source["broj_zadataka"];
	        this.broj_clanova = source["broj_clanova"];
	    }
	}

}

export namespace services {
	
	export class LoginResponse {
	    user?: models.Korisnici;
	    success: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new LoginResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.user = this.convertValues(source["user"], models.Korisnici);
	        this.success = source["success"];
	        this.message = source["message"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

