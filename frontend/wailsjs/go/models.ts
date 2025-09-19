export namespace models {
	
	export class Dokumenti {
	    dokument_id: number;
	    projekat_id?: number;
	    naziv_dokumenta: string;
	    folder_id?: number;
	    opis?: string;
	    tip_dokumenta?: string;
	    jezik_dokumenta?: string;
	    radni_tok_id?: number;
	    trenutna_faza_id?: number;
	    kreirao_korisnik_id: number;
	    // Go type: time
	    datuma_postavke: any;
	    // Go type: time
	    poslednja_izmena?: any;
	    naziv_projekta?: string;
	    ime_kreirao?: string;
	    naziv_faze?: string;
	    broj_verzija?: number;
	
	    static createFrom(source: any = {}) {
	        return new Dokumenti(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dokument_id = source["dokument_id"];
	        this.projekat_id = source["projekat_id"];
	        this.naziv_dokumenta = source["naziv_dokumenta"];
	        this.folder_id = source["folder_id"];
	        this.opis = source["opis"];
	        this.tip_dokumenta = source["tip_dokumenta"];
	        this.jezik_dokumenta = source["jezik_dokumenta"];
	        this.radni_tok_id = source["radni_tok_id"];
	        this.trenutna_faza_id = source["trenutna_faza_id"];
	        this.kreirao_korisnik_id = source["kreirao_korisnik_id"];
	        this.datuma_postavke = this.convertValues(source["datuma_postavke"], null);
	        this.poslednja_izmena = this.convertValues(source["poslednja_izmena"], null);
	        this.naziv_projekta = source["naziv_projekta"];
	        this.ime_kreirao = source["ime_kreirao"];
	        this.naziv_faze = source["naziv_faze"];
	        this.broj_verzija = source["broj_verzija"];
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
	export class Tagovi {
	    tag_id: number;
	    naziv_taga: string;
	
	    static createFrom(source: any = {}) {
	        return new Tagovi(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tag_id = source["tag_id"];
	        this.naziv_taga = source["naziv_taga"];
	    }
	}
	export class UploadDocumentRequest {
	    naziv_dokumenta: string;
	    projekat_id?: number;
	    folder_id?: number;
	    opis: string;
	    tip_dokumenta: string;
	    jezik_dokumenta: string;
	    tagovi: string[];
	
	    static createFrom(source: any = {}) {
	        return new UploadDocumentRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.naziv_dokumenta = source["naziv_dokumenta"];
	        this.projekat_id = source["projekat_id"];
	        this.folder_id = source["folder_id"];
	        this.opis = source["opis"];
	        this.tip_dokumenta = source["tip_dokumenta"];
	        this.jezik_dokumenta = source["jezik_dokumenta"];
	        this.tagovi = source["tagovi"];
	    }
	}
	export class VerzijeDokumenata {
	    verzija_id: number;
	    dokument_id: number;
	    verzija_oznaka?: string;
	    putanja_do_fajla: string;
	    velicina_fajla_mb?: number;
	    postavio_korisnik_id: number;
	    // Go type: time
	    datuma_postavke: any;
	
	    static createFrom(source: any = {}) {
	        return new VerzijeDokumenata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.verzija_id = source["verzija_id"];
	        this.dokument_id = source["dokument_id"];
	        this.verzija_oznaka = source["verzija_oznaka"];
	        this.putanja_do_fajla = source["putanja_do_fajla"];
	        this.velicina_fajla_mb = source["velicina_fajla_mb"];
	        this.postavio_korisnik_id = source["postavio_korisnik_id"];
	        this.datuma_postavke = this.convertValues(source["datuma_postavke"], null);
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

